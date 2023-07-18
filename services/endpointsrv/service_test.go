package endpointsrv

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

func TestStepStatusAndConclusion(t *testing.T) {
	builder := newDBCheckStepBuilder()
	stepWithModelOutput := builder.withModelOutput("foo")

	t.Run("pending", func(t *testing.T) {
		step := builder.mustBuild()

		status, conclusion := stepStatusAndConclusion(step)

		require.Equal(t, types.CheckPending, status)
		require.Nil(t, conclusion)
	})

	t.Run("in progress", func(t *testing.T) {
		t.Run("no assertion", func(t *testing.T) {
			step := stepWithModelOutput.mustBuild()

			status, conclusion := stepStatusAndConclusion(step)

			require.Equal(t, types.CheckInProgress, status)
			require.Nil(t, conclusion)
		})

		t.Run("with empty assertion", func(t *testing.T) {
			step := stepWithModelOutput.withAssertionOutput("").mustBuild()

			status, conclusion := stepStatusAndConclusion(step)

			require.Equal(t, types.CheckInProgress, status)
			require.Nil(t, conclusion)
		})
	})

	t.Run("completed", func(t *testing.T) {
		type testCase struct {
			name               string
			step               func() db.UnweaveEndpointCheckStep
			expectedConclusion types.CheckConclusion
		}

		testCases := []testCase{
			{
				name: "success",
				step: func() db.UnweaveEndpointCheckStep {
					return stepWithModelOutput.withAssertionOutput("success").mustBuild()
				},
				expectedConclusion: types.CheckSuccess,
			},
			{
				name: "failure",
				step: func() db.UnweaveEndpointCheckStep {
					return stepWithModelOutput.withAssertionOutput("failure").mustBuild()
				},
				expectedConclusion: types.CheckFailure,
			},
			{
				name: "error",
				step: func() db.UnweaveEndpointCheckStep {
					return stepWithModelOutput.withAssertionOutput("foo").mustBuild()
				},
				expectedConclusion: types.CheckError,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				step := tc.step()

				status, conclusion := stepStatusAndConclusion(step)

				require.Equal(t, types.CheckCompleted, status)
				require.NotNil(t, conclusion)
				require.Equal(t, tc.expectedConclusion, *conclusion)
			})
		}
	})
}

func TestCheckStatusAndConclusion(t *testing.T) {
	t.Run("in progress", func(t *testing.T) {
		type testCase struct {
			name  string
			steps []types.EndpointCheckStep
		}

		testCases := []testCase{
			{
				name: "all pending",
				steps: []types.EndpointCheckStep{
					{Status: types.CheckPending},
					{Status: types.CheckPending},
					{Status: types.CheckPending},
				},
			},
			{
				name: "some in progress",
				steps: []types.EndpointCheckStep{
					{Status: types.CheckPending},
					{Status: types.CheckInProgress},
					{Status: types.CheckPending},
				},
			},
			{
				name: "some completed",
				steps: []types.EndpointCheckStep{
					{Status: types.CheckPending},
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				status, conclusion := checkStatusAndConclusion(tc.steps)

				require.Equal(t, types.CheckInProgress, status)
				require.Nil(t, conclusion)
			})
		}
	})

	t.Run("completed", func(t *testing.T) {
		type testCase struct {
			name               string
			steps              []types.EndpointCheckStep
			expectedConclusion types.CheckConclusion
		}

		testCases := []testCase{
			{
				name: "all success",
				steps: []types.EndpointCheckStep{
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
				},
				expectedConclusion: types.CheckSuccess,
			},
			{
				name: "some failure",
				steps: []types.EndpointCheckStep{
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
					{Status: types.CheckCompleted, Conclusion: &types.CheckFailure},
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
				},
				expectedConclusion: types.CheckFailure,
			},
			{
				name: "some error",
				steps: []types.EndpointCheckStep{
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
					{Status: types.CheckCompleted, Conclusion: &types.CheckError},
					{Status: types.CheckCompleted, Conclusion: &types.CheckSuccess},
				},
				expectedConclusion: types.CheckError,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				status, conclusion := checkStatusAndConclusion(tc.steps)

				require.Equal(t, types.CheckCompleted, status)
				require.NotNil(t, conclusion)
				require.Equal(t, tc.expectedConclusion, *conclusion)
			})
		}
	})
}

type dbEndpointCheckStepBuilder struct {
	object db.UnweaveEndpointCheckStep
}

func newDBCheckStepBuilder() dbEndpointCheckStepBuilder {
	return dbEndpointCheckStepBuilder{}
}

func (b dbEndpointCheckStepBuilder) withModelOutput(o string) dbEndpointCheckStepBuilder {
	b.object.Output = sql.NullString{String: o, Valid: true}
	return b
}

func (b dbEndpointCheckStepBuilder) withAssertionOutput(o string) dbEndpointCheckStepBuilder {
	b.object.Assertion = sql.NullString{String: o, Valid: true}
	return b
}

func (b dbEndpointCheckStepBuilder) build() (db.UnweaveEndpointCheckStep, error) {
	return b.object, nil
}

func (b dbEndpointCheckStepBuilder) mustBuild() db.UnweaveEndpointCheckStep {
	step, err := b.build()
	if err != nil {
		panic(err)
	}

	return step
}
