package domain_test

import (
	"errors"
	"testing"

	adaptermocks "github.com/mouse-blink/gooze/internal/adapter/mocks"
	controllermocks "github.com/mouse-blink/gooze/internal/controller/mocks"
	domain "github.com/mouse-blink/gooze/internal/domain"
	domainmocks "github.com/mouse-blink/gooze/internal/domain/mocks"
	m "github.com/mouse-blink/gooze/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWorkflow_Test_Success(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	sources := []m.Source{
		{
			Origin: &m.File{FullPath: "test.go", Hash: "hash1"},
			Test:   &m.File{FullPath: "test_test.go", Hash: "test_hash1"},
		},
	}

	mutations := []m.Mutation{
		{ID: 1, Source: sources[0], Type: m.MutationArithmetic},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mock.Anything, mock.Anything).Return().Once()
	mockUI.EXPECT().DisplayCompletedTestInfo(mock.Anything, mock.Anything).Return().Once()
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	mockOrchestrator.EXPECT().TestMutation(mock.Anything).Return(m.Result{}, nil)
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.Anything).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockFSAdapter.AssertExpectations(t)
	mockMutagen.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
	mockOrchestrator.AssertExpectations(t)
}

func TestWorkflow_Test_GetSourcesError(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	testErr := errors.New("failed to get sources")
	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(nil, testErr)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports: "reports.json",
	}
	err := wf.Test(args)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get sources")
}

func TestWorkflow_Test_GenerateMutationsError(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	sources := []m.Source{
		{Origin: &m.File{FullPath: "test.go", Hash: "hash1"}},
	}

	testErr := errors.New("failed to generate mutations")
	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(nil, testErr)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports: "reports.json",
	}
	err := wf.Test(args)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "generate mutations")
}

func TestWorkflow_Test_TestMutationError(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	sources := []m.Source{
		{Origin: &m.File{FullPath: "test.go", Hash: "hash1"}},
	}

	mutations := []m.Mutation{
		{ID: 1, Source: sources[0]},
	}

	testErr := errors.New("failed to test mutation")
	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mock.Anything, mock.Anything).Return().Once()
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	mockOrchestrator.EXPECT().TestMutation(mock.Anything).Return(nil, testErr)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "errors occurred during mutation testing")
}

func TestWorkflow_Test_SaveReportsError(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	sources := []m.Source{
		{Origin: &m.File{FullPath: "test.go", Hash: "hash1"}},
	}

	mutations := []m.Mutation{
		{ID: 1, Source: sources[0]},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mock.Anything, mock.Anything).Return().Once()
	mockUI.EXPECT().DisplayCompletedTestInfo(mock.Anything, mock.Anything).Return().Once()
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	mockOrchestrator.EXPECT().TestMutation(mock.Anything).Return(m.Result{}, nil)

	saveErr := errors.New("failed to save reports")
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.Anything).Return(saveErr)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save reports")
}

func TestWorkflow_Test_NoMutations(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	sources := []m.Source{
		{Origin: &m.File{FullPath: "test.go", Hash: "hash1"}},
	}

	// No mutations generated
	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return([]m.Mutation{}, nil)
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.MatchedBy(func(reports []m.Report) bool {
		return len(reports) == 0
	})).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockFSAdapter.AssertExpectations(t)
	mockMutagen.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
}

func TestWorkflow_Test_MultipleThreads(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	source := m.Source{
		Origin: &m.File{FullPath: "test.go", Hash: "hash1"},
	}
	sources := []m.Source{source}

	mutations := []m.Mutation{
		{ID: 0, Source: source},
		{ID: 1, Source: source},
		{ID: 2, Source: source},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mock.Anything, mock.Anything).Return().Times(3)
	mockUI.EXPECT().DisplayCompletedTestInfo(mock.Anything, mock.Anything).Return().Times(3)
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	mockOrchestrator.EXPECT().TestMutation(mock.Anything).Return(m.Result{}, nil).Times(3)
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.MatchedBy(func(reports []m.Report) bool {
		return len(reports) == 3
	})).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports:         "reports.json",
		Threads:         4,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockOrchestrator.AssertExpectations(t)
}

func TestWorkflow_Test_WithSharding(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	source := m.Source{
		Origin: &m.File{FullPath: "test.go", Hash: "hash1"},
	}

	// 6 mutations total
	mutations := []m.Mutation{
		{ID: 0, Source: source},
		{ID: 1, Source: source},
		{ID: 2, Source: source},
		{ID: 3, Source: source},
		{ID: 4, Source: source},
		{ID: 5, Source: source},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mock.Anything, mock.Anything).Return().Times(2)
	mockUI.EXPECT().DisplayCompletedTestInfo(mock.Anything, mock.Anything).Return().Times(2)
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return([]m.Source{source}, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	// Only 2 mutations should be tested (IDs 0 and 3, since shardIndex=0, totalShards=3)
	mockOrchestrator.EXPECT().TestMutation(mock.Anything).Return(m.Result{}, nil).Times(2)
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.MatchedBy(func(reports []m.Report) bool {
		return len(reports) == 2
	})).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		EstimateArgs: domain.EstimateArgs{
			Paths: []m.Path{"test.go"},
		},
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 3,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockOrchestrator.AssertExpectations(t)
}

func TestWorkflow_Test_MultipleSources(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	source1 := m.Source{
		Origin: &m.File{FullPath: "file1.go", Hash: "hash1"},
	}
	source2 := m.Source{
		Origin: &m.File{FullPath: "file2.go", Hash: "hash2"},
	}

	mutations1 := []m.Mutation{
		{ID: 0, Source: source1},
		{ID: 1, Source: source1},
	}
	mutations2 := []m.Mutation{
		{ID: 2, Source: source2},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(mock.Anything).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mock.Anything, mock.Anything).Return().Times(3)
	mockUI.EXPECT().DisplayCompletedTestInfo(mock.Anything, mock.Anything).Return().Times(3)
	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return([]m.Source{source1, source2}, nil)
	mockMutagen.EXPECT().GenerateMutation(source1, 0, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations1, nil)
	mockMutagen.EXPECT().GenerateMutation(source2, 2, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations2, nil)
	mockOrchestrator.EXPECT().TestMutation(mock.Anything).Return(m.Result{}, nil).Times(3)
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.MatchedBy(func(reports []m.Report) bool {
		return len(reports) == 3
	})).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{

		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockFSAdapter.AssertExpectations(t)
	mockMutagen.AssertExpectations(t)
	mockOrchestrator.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
}

func TestWorkflow_NewWorkflowV2(t *testing.T) {
	// Arrange
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	// Act
	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Assert
	require.NotNil(t, wf)
	assert.Implements(t, (*domain.Workflow)(nil), wf)
}

func TestWorkflow_TestWithSurvivedMutation(t *testing.T) {
	// Arrange - This test specifically checks that survived mutations include diff data
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	diffCode := []byte("--- original\n+++ mutated\n@@ -1,1 +1,1 @@\n-\treturn 3 + 5\n+\treturn 3 - 5\n")

	sources := []m.Source{
		{
			Origin: &m.File{FullPath: "test.go", Hash: "hash1"},
			Test:   &m.File{FullPath: "test_test.go", Hash: "test_hash1"},
		},
	}

	mutations := []m.Mutation{
		{
			ID:       0,
			Source:   sources[0],
			Type:     m.MutationArithmetic,
			DiffCode: diffCode,
		},
	}

	// Mock a survived mutation result
	survivedResult := m.Result{
		m.MutationArithmetic: []struct {
			MutationID string
			Status     m.TestStatus
			Err        error
		}{{MutationID: "0", Status: m.Survived}},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(1).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mutations[0], 0).Return().Once()
	mockUI.EXPECT().DisplayCompletedTestInfo(mutations[0], survivedResult).Return().Once()

	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	mockOrchestrator.EXPECT().TestMutation(mutations[0]).Return(survivedResult, nil)

	// Verify that the report includes the diff for survived mutations
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.MatchedBy(func(reports []m.Report) bool {
		if len(reports) != 1 {
			return false
		}
		report := reports[0]
		// Check that diff is included for survived mutation
		return report.Diff != nil && string(*report.Diff) == string(diffCode)
	})).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockFSAdapter.AssertExpectations(t)
	mockMutagen.AssertExpectations(t)
	mockOrchestrator.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
	mockUI.AssertExpectations(t)
}

func TestWorkflow_TestWithKilledMutation(t *testing.T) {
	// Arrange - This test checks that killed mutations do NOT include diff data
	mockFSAdapter := new(adaptermocks.MockSourceFSAdapter)
	mockReportStore := new(adaptermocks.MockReportStore)
	mockUI := new(controllermocks.MockUI)
	mockOrchestrator := new(domainmocks.MockOrchestrator)
	mockMutagen := new(domainmocks.MockMutagen)

	diffCode := []byte("--- original\n+++ mutated\n@@ -1,1 +1,1 @@\n-\treturn 3 + 5\n+\treturn 3 - 5\n")

	sources := []m.Source{
		{
			Origin: &m.File{FullPath: "test.go", Hash: "hash1"},
			Test:   &m.File{FullPath: "test_test.go", Hash: "test_hash1"},
		},
	}

	mutations := []m.Mutation{
		{
			ID:       0,
			Source:   sources[0],
			Type:     m.MutationArithmetic,
			DiffCode: diffCode,
		},
	}

	// Mock a killed mutation result
	killedResult := m.Result{
		m.MutationArithmetic: []struct {
			MutationID string
			Status     m.TestStatus
			Err        error
		}{{MutationID: "0", Status: m.Killed}},
	}

	mockUI.EXPECT().Start(mock.Anything).Return(nil).Once()
	mockUI.EXPECT().Wait().Return().Once()
	mockUI.EXPECT().Close().Return().Once()
	mockUI.EXPECT().DisplayConcurencyInfo(mock.Anything, mock.Anything, mock.Anything).Return()
	mockUI.EXPECT().DusplayUpcomingTestsInfo(1).Return()
	mockUI.EXPECT().DisplayStartingTestInfo(mutations[0], 0).Return().Once()
	mockUI.EXPECT().DisplayCompletedTestInfo(mutations[0], killedResult).Return().Once()

	mockFSAdapter.EXPECT().Get(mock.Anything, domain.DefaultIgnorePattern).Return(sources, nil)
	mockMutagen.EXPECT().GenerateMutation(mock.Anything, mock.Anything, domain.DefaultMutations[0], domain.DefaultMutations[1], domain.DefaultMutations[2], domain.DefaultMutations[3], domain.DefaultMutations[4]).Return(mutations, nil)
	mockOrchestrator.EXPECT().TestMutation(mutations[0]).Return(killedResult, nil)

	// Verify that the report does NOT include diff for killed mutations
	mockReportStore.EXPECT().SaveReports(mock.Anything, mock.MatchedBy(func(reports []m.Report) bool {
		if len(reports) != 1 {
			return false
		}
		report := reports[0]
		// Check that diff is NOT included for killed mutation
		return report.Diff == nil
	})).Return(nil)

	wf := domain.NewWorkflow(mockFSAdapter, mockReportStore, mockUI, mockOrchestrator, mockMutagen)

	// Act
	args := domain.TestArgs{
		Reports:         "reports.json",
		Threads:         1,
		ShardIndex:      0,
		TotalShardCount: 1,
	}
	err := wf.Test(args)

	// Assert
	assert.NoError(t, err)
	mockFSAdapter.AssertExpectations(t)
	mockMutagen.AssertExpectations(t)
	mockOrchestrator.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
	mockUI.AssertExpectations(t)
}
