//nolint:noctx,godox
package endpointsrv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/services/evalsrv"
	"github.com/unweave/unweave/services/execsrv"
	"github.com/unweave/unweave/tools/random"
	"go.jetpack.io/typeid"
)

type Driver interface {
	EndpointDriverName() string
	EndpointProvider() types.Provider

	EndpointCreate(
		ctx context.Context,
		project,
		endpointID,
		subdomain string) (string, error)

	EndpointVersionCreate(
		ctx context.Context,
		project,
		endpointID,
		versionID,
		execID string,
		internalPort int32) (string, error)

	EndpointVersionPromote(
		ctx context.Context,
		endpointID,
		versionID string,
		internalPort int32) error
}

type Service interface {
	EndpointExecCreate(ctx context.Context, projectID, execID, name string) (types.Endpoint, error)
	EndpointGet(ctx context.Context, projectID, endpointID string) (types.Endpoint, error)
	EndpointList(ctx context.Context, projectID string) ([]types.EndpointListItem, error)

	RunEndpointEvals(ctx context.Context, projectID, endpointID string) (string, error)
	EndpointAttachEval(ctx context.Context, endpointID, evalID string) error
	EndpointCheckStatus(ctx context.Context, checkID string) (types.EndpointCheck, error)

	EndpointVersionCreate(ctx context.Context, projectID, parentEndpointID, execID string, promote bool) (types.EndpointVersion, error)
}

type EndpointService struct {
	store  Store
	evals  evalsrv.Service
	execs  execsrv.Service
	driver Driver
}

var _ Service = (*EndpointService)(nil)

type Store interface {
	EndpointCreate(ctx context.Context, arg db.EndpointCreateParams) error
	EndpointGet(ctx context.Context, arg db.EndpointGetParams) (db.UnweaveEndpoint, error)
	EndpointsForProject(ctx context.Context, id string) ([]db.UnweaveEndpoint, error)
	EndpointEval(ctx context.Context, endpointID string) ([]db.UnweaveEndpointEval, error)
	EndpointEvalAttach(ctx context.Context, arg db.EndpointEvalAttachParams) error
	EndpointCheckCreate(ctx context.Context, arg db.EndpointCheckCreateParams) error
	EndpointCheckStepCreate(ctx context.Context, arg db.EndpointCheckStepCreateParams) error
	EndpointCheckStepUpdate(ctx context.Context, arg db.EndpointCheckStepUpdateParams) error
	EndpointCheckSteps(ctx context.Context, checkID string) ([]db.UnweaveEndpointCheckStep, error)
	EndpointCheck(ctx context.Context, checkID string) (db.UnweaveEndpointCheck, error)

	EndpointVersion(ctx context.Context, versionID string) (db.UnweaveEndpointVersion, error)
	EndpointVersionCreate(ctx context.Context, arg db.EndpointVersionCreateParams) error
	EndpointVersionList(ctx context.Context, endpointID string) ([]db.UnweaveEndpointVersion, error)
	EndpointVersionDemote(ctx context.Context, endpointID string) error
	EndpointVersionPromote(ctx context.Context, id string) error

	Tx(txFunc func(db.Querier) error) error
}

func NewEndpointService(
	store Store,
	evals evalsrv.Service,
	execs execsrv.Service,
	driver Driver,
) *EndpointService {
	return &EndpointService{
		store:  store,
		evals:  evals,
		execs:  execs,
		driver: driver,
	}
}

// EndpointExecCreate creates a new endpoint (with http address)
// it also creates the first endpoint version for the exec
// promotes that version to primary, and serves traffic to it.
func (e *EndpointService) EndpointExecCreate(ctx context.Context, projectID, execID, endpointName string) (types.Endpoint, error) {
	endpointID := typeid.Must(typeid.New("end")).String()

	if endpointName == "" {
		endpointName = random.GenerateRandomPhrase(3, "-")
	}

	end, err := e.createEndpoint(ctx, endpointID, endpointName, projectID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("create endpoint: %w", err)
	}

	return end, nil
}

func (e *EndpointService) createEndpoint(
	ctx context.Context,
	endpointID,
	endpointName,
	projectID string,
) (types.Endpoint, error) {
	subdomain := fmt.Sprint(endpointName, "-", random.GenerateRandomLower(5))

	httpAddr, err := e.driver.EndpointCreate(
		ctx,
		projectID,
		endpointID,
		subdomain,
	)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("drvier create endpoint: %w", err)
	}

	end := types.Endpoint{
		ID:          endpointID,
		Name:        endpointName,
		Icon:        "🚀",
		ProjectID:   projectID,
		HTTPAddress: httpAddr,
		EvalIDs:     []string{},
		Status:      "",
		Versions:    []types.EndpointVersion{},
		CreatedAt:   time.Now(),
	}

	dbe := db.EndpointCreateParams{
		ID:          end.ID,
		Name:        endpointName,
		Icon:        "🚀",
		ProjectID:   end.ProjectID,
		HttpAddress: end.HTTPAddress,
		CreatedAt:   end.CreatedAt,
	}

	if err := e.store.EndpointCreate(ctx, dbe); err != nil {
		return types.Endpoint{}, fmt.Errorf("save endpoint: %w", err)
	}

	return end, nil
}

func (e *EndpointService) EndpointList(ctx context.Context, projectID string) ([]types.EndpointListItem, error) {
	dbEnds, err := e.store.EndpointsForProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get endpoint: %w", err)
	}

	out := []types.EndpointListItem{}

	for _, end := range dbEnds {
		out = append(out, types.EndpointListItem{
			ID:          end.ID,
			Name:        end.Name,
			Icon:        end.Icon,
			ProjectID:   end.ProjectID,
			HTTPAddress: end.HttpAddress,
			CreatedAt:   end.CreatedAt,
		})
	}

	return out, nil
}

func (e *EndpointService) EndpointGet(ctx context.Context, projectID, endpointID string) (types.Endpoint, error) {
	arg := db.EndpointGetParams{
		ID:        endpointID,
		ProjectID: projectID,
	}

	end, err := e.store.EndpointGet(ctx, arg)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get endpoint: %w", err)
	}

	endEvals, err := e.store.EndpointEval(ctx, end.ID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("get endpoint evals: %w", err)
	}

	ids := make([]string, len(endEvals))
	for i, eval := range endEvals {
		ids[i] = eval.EvalID
	}

	versions, err := e.endpointVersions(ctx, endpointID)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("versions: %w", err)
	}

	endpoint := types.Endpoint{
		ID:          end.ID,
		Name:        end.Name,
		Icon:        end.Icon,
		ProjectID:   end.ProjectID,
		HTTPAddress: end.HttpAddress,
		EvalIDs:     ids,
		Status:      "",
		Versions:    versions,
		CreatedAt:   end.CreatedAt,
	}

	return endpoint, nil
}

func (e *EndpointService) endpointVersions(ctx context.Context, endpointID string) ([]types.EndpointVersion, error) {
	vers, err := e.store.EndpointVersionList(ctx, endpointID)
	if err != nil {
		return nil, fmt.Errorf("version list: %w", err)
	}

	out := make([]types.EndpointVersion, len(vers))

	for idx, ver := range vers {
		out[idx] = types.EndpointVersion{
			ID:          ver.ID,
			ExecID:      ver.ExecID,
			HTTPAddress: ver.HttpAddress,
			Status:      "",
			Primary:     ver.PrimaryVersion,
			CreatedAt:   ver.CreatedAt,
		}
	}

	return out, nil
}

func (e *EndpointService) EndpointAttachEval(ctx context.Context, endpointID, evalID string) error {
	if err := e.store.EndpointEvalAttach(ctx, db.EndpointEvalAttachParams{
		EndpointID: endpointID,
		EvalID:     evalID,
	}); err != nil {
		return fmt.Errorf("attach eval: %w", err)
	}

	return nil
}

func (e *EndpointService) RunEndpointEvals(ctx context.Context, projectID, endpointID string) (string, error) {
	checkID := typeid.Must(typeid.New("check")).String()

	endpoint, err := e.EndpointGet(ctx, projectID, endpointID)
	if err != nil {
		return "", fmt.Errorf("get endpoint: %w", err)
	}

	evals, err := e.evals.Evals(ctx, endpoint.EvalIDs)
	if err != nil {
		return "", fmt.Errorf("get evals: %w", err)
	}

	if err := verifyCanRunChecks(endpoint, evals); err != nil {
		return "", fmt.Errorf("verify checks: %w", err)
	}

	if err := e.store.EndpointCheckCreate(ctx, db.EndpointCheckCreateParams{
		ID:         checkID,
		EndpointID: endpoint.ID,
		ProjectID:  endpoint.ProjectID,
	}); err != nil {
		return "", fmt.Errorf("create eval check: %w", err)
	}

	checker, err := newEndpointChecker(ctx, checkID)
	if err != nil {
		return "", fmt.Errorf("new endpoint checker: %w", err)
	}

	err = checker.CreateCheckSteps(ctx, e.store, endpoint, evals)
	if err != nil {
		return "", fmt.Errorf("create endpoint checks: %w", err)
	}

	return checkID, checker.Run(context.Background())
}

func verifyCanRunChecks(endpoint types.Endpoint, evals []types.Eval) error {
	if endpoint.HTTPAddress == "" {
		return errors.New("endpoint must be exposed and have hostname attached")
	}

	if !allHaveHTTPServiceHostname(evals...) {
		return errors.New("evals must have http service exposed")
	}

	if len(evals) == 0 {
		return errors.New("no evals on endpoint")
	}

	return nil
}

func allHaveHTTPServiceHostname(evals ...types.Eval) bool {
	for _, e := range evals {
		if e.HTTPEndpoint == "" {
			return false
		}
	}

	return true
}

func (e *EndpointService) EndpointCheckStatus(ctx context.Context, checkID string) (types.EndpointCheck, error) {
	steps, err := e.store.EndpointCheckSteps(ctx, checkID)
	if err != nil {
		return types.EndpointCheck{}, fmt.Errorf("get steps: %w", err)
	}

	out := make([]types.EndpointCheckStep, len(steps))
	for idx, step := range steps {
		status, conclusion := stepStatusAndConclusion(step)
		out[idx] = types.EndpointCheckStep{
			StepID:     step.ID,
			EvalID:     step.EvalID,
			Input:      []byte(step.Input.String),
			Output:     []byte(step.Output.String),
			Assertion:  step.Assertion.String,
			Status:     status,
			Conclusion: conclusion,
		}
	}

	status, conclusion := checkStatusAndConclusion(out)
	return types.EndpointCheck{
		CheckID:    checkID,
		Steps:      out,
		Status:     status,
		Conclusion: conclusion,
	}, nil
}

func (e *EndpointService) EndpointVersionCreate(
	ctx context.Context,
	projectID,
	endpointID,
	execID string,
	promote bool,
) (types.EndpointVersion, error) {
	end, err := e.EndpointGet(ctx, projectID, endpointID)
	if err != nil {
		return types.EndpointVersion{}, fmt.Errorf("endpoint get: %w", err)
	}

	exec, err := e.execs.Get(ctx, execID)
	if err != nil {
		return types.EndpointVersion{}, fmt.Errorf("exec get: %w", err)
	}

	if exec.Network.HTTPService == nil {
		return types.EndpointVersion{}, &types.Error{
			Code:       400,
			Message:    "Cannot create endpoint version on exec with no port",
			Suggestion: "Create the exec with a port",
			Provider:   e.driver.EndpointProvider(),
		}
	}

	return e.createAttachEndpointVersion(ctx, end, exec.ID, exec.Network.HTTPService.InternalPort, promote)
}

func (e *EndpointService) createAttachEndpointVersion(
	ctx context.Context,
	end types.Endpoint,
	execID string,
	internalPort int32,
	promote bool,
) (types.EndpointVersion, error) {
	versionID := typeid.Must(typeid.New("version")).String()

	log.Debug().
		Str(types.ProjectIDCtxKey, end.ProjectID).
		Str(types.EndpointIDCtxKey, end.ID).
		Str(types.ExecIDCtxKey, execID).
		Str(types.VersionIDCtxKey, versionID).
		Bool("promote", promote).
		Msg("creating endpoint version")

	httpAddr, err := e.driver.EndpointVersionCreate(
		ctx,
		end.ProjectID,
		end.ID,
		versionID,
		execID,
		internalPort,
	)
	if err != nil {
		return types.EndpointVersion{}, fmt.Errorf("version create: %w", err)
	}

	args := db.EndpointVersionCreateParams{
		ID:          versionID,
		EndpointID:  end.ID,
		ExecID:      execID,
		ProjectID:   end.ProjectID,
		HttpAddress: httpAddr,
		CreatedAt:   time.Now(),
	}

	if err := e.store.EndpointVersionCreate(ctx, args); err != nil {
		return types.EndpointVersion{}, fmt.Errorf("version store: %w", err)
	}

	version := types.EndpointVersion{
		ID:          versionID,
		ExecID:      execID,
		HTTPAddress: httpAddr,
		Primary:     promote,
		CreatedAt:   args.CreatedAt,
	}

	if version.Primary {
		if err := e.setPrimary(ctx, end, version, internalPort); err != nil {
			return types.EndpointVersion{}, fmt.Errorf("promote: %w", err)
		}
	}

	return version, nil
}

func (e *EndpointService) setPrimary(
	ctx context.Context,
	end types.Endpoint,
	version types.EndpointVersion,
	internalPort int32,
) error {
	demoteVersions(&end)

	if err := e.driver.EndpointVersionPromote(
		ctx,
		end.ID,
		version.ID,
		internalPort,
	); err != nil {
		return fmt.Errorf("update version: %w", err)
	}

	transaction := func(q db.Querier) error {
		if err := q.EndpointVersionDemote(ctx, end.ID); err != nil {
			return fmt.Errorf("demote: %w", err)
		}

		if err := q.EndpointVersionPromote(ctx, version.ID); err != nil {
			return fmt.Errorf("promote: %w", err)
		}

		return nil
	}

	if err := e.store.Tx(transaction); err != nil {
		return fmt.Errorf("update primary: %w", err)
	}

	return nil
}

func demoteVersions(end *types.Endpoint) {
	for i := range end.Versions {
		end.Versions[i].Primary = false
	}
}

// stepStatusAndConclusion infers the status of a step based on the contents of the database.
func stepStatusAndConclusion(step db.UnweaveEndpointCheckStep) (types.CheckStatus, *types.CheckConclusion) {
	hasModelOutput := step.Output.Valid
	hasAssertionOutput := step.Assertion.Valid

	if !hasModelOutput {
		return types.CheckPending, nil
	}

	if !hasAssertionOutput {
		return types.CheckInProgress, nil
	}

	switch step.Assertion.String {
	case "":
		return types.CheckInProgress, nil
	case "success":
		return types.CheckCompleted, &types.CheckSuccess
	case "failure":
		return types.CheckCompleted, &types.CheckFailure
	default:
		return types.CheckCompleted, &types.CheckError
	}
}

// checkStatusAndConclusion infers the status of a check based on the status of the steps.
func checkStatusAndConclusion(steps []types.EndpointCheckStep) (types.CheckStatus, *types.CheckConclusion) {
	conclusionsHave := map[types.CheckConclusion]bool{}

	for _, step := range steps {
		if step.Status != types.CheckCompleted {
			return types.CheckInProgress, nil
		}

		conclusionsHave[*step.Conclusion] = true
	}

	if conclusionsHave[types.CheckError] {
		return types.CheckCompleted, &types.CheckError
	}

	if conclusionsHave[types.CheckFailure] {
		return types.CheckCompleted, &types.CheckFailure
	}

	return types.CheckCompleted, &types.CheckSuccess
}
