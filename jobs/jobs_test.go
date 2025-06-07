package jobs

import (
	"reflect"
	"sort"
	"testing"

	"github.com/arunjeyaprasad/golive/models"
	"github.com/arunjeyaprasad/golive/streamer"
)

func setup() {
	// Clear the jobs map before each test
	jobs = make(map[string]models.Job)
	// Clear the jobProcessMap before each test
	jobProcessMap = make(map[string]*streamer.StreamingProcess)
}

func TestStopJob(t *testing.T) {
	type args struct {
		jobID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Invalid job id",
			args: args{
				jobID: "invalid-job-id",
			},
			wantErr: false,
		},
	}
	setup() // Call setup to clear the jobs map before each test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StopJob(tt.args.jobID); (err != nil) != tt.wantErr {
				t.Errorf("StopJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateJob(t *testing.T) {
	type args struct {
		request models.JobCreateRequest
	}
	tests := []struct {
		name string
		args args
		want *models.Job
	}{
		{
			name: "Create job with valid description",
			args: args{
				request: models.JobCreateRequest{
					Description: "Test job",
				},
			},
		},
	}
	setup() // Call setup to clear the jobs map before each test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateJob(tt.args.request)
			if got.ID == "" {
				t.Errorf("CreateJob() returned job with empty ID")
			}
			if got.Configuration.Description != tt.args.request.Description {
				t.Errorf("CreateJob() returned job with description = %v, want %v", got.Configuration.Description, tt.args.request.Description)
			}
			if got.Status != string(JobStatusCreated) {
				t.Errorf("CreateJob() returned job with status = %v, want %v", got.Status, JobStatusCreated)
			}
			if got.CreatedAt == "" {
				t.Errorf("CreateJob() returned job with empty CreatedAt")
			}
			// Check if the job is added to the jobs map
			if len(jobs) == 0 {
				t.Errorf("CreateJob() did not add job to jobs map")
			}
		})
	}
}

func TestGetJobs(t *testing.T) {
	// Initialize the jobs map with some test data
	testJobs := []models.Job{
		{
			ID:        "job1",
			Status:    string(JobStatusCreated),
			CreatedAt: "2023-10-01T00:00:00Z",
			Configuration: models.JobCreateRequest{
				Description: "Test job 1",
			},
		},
		{
			ID:        "job2",
			Status:    string(JobStatusRunning),
			CreatedAt: "2023-10-02T00:00:00Z",
			Configuration: models.JobCreateRequest{
				Description: "Test job 2",
			},
		},
		{
			ID:        "job3",
			Status:    string(JobStatusCompleted),
			CreatedAt: "2023-10-03T00:00:00Z",
			Configuration: models.JobCreateRequest{
				Description: "Test job 3",
			},
		},
	}
	tests := []struct {
		name string
		want []models.Job
	}{
		{
			name: "Get all jobs",
			want: testJobs,
		},
	}
	setup() // Call setup to clear the jobs map before each test

	for _, job := range testJobs {
		jobs[job.ID] = job
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetJobs()
			want := tt.want
			if len(got) != len(want) {
				t.Errorf("GetJobs() length = %d, want %d", len(got), len(want))
			}
			// Check if the jobs match
			sort.Slice(got, func(i, j int) bool {
				return got[i].ID < got[j].ID
			})
			sort.Slice(want, func(i, j int) bool {
				return want[i].ID < want[j].ID
			})

			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetJobs() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name  string
		args  args
		want  *models.Job
		want1 bool
	}{
		{
			name: "Get existing job",
			args: args{
				id: "job1",
			},
			want: &models.Job{
				ID:        "job1",
				Status:    string(JobStatusCreated),
				CreatedAt: "2023-10-01T00:00:00Z",
				Configuration: models.JobCreateRequest{
					Description: "Test job 1",
				},
			},
			want1: true,
		},
		{
			name: "Get non-existing job",
			args: args{
				id: "non-existing-job",
			},
			want:  nil,
			want1: false,
		},
	}
	setup() // Call setup to clear the jobs map before each test
	// Add a test job to the jobs map
	j := tests[0].want
	jobs[j.ID] = *j
	// Add a second job to the jobs map
	jobs["job2"] = models.Job{
		ID:        "job2",
		Status:    string(JobStatusRunning),
		CreatedAt: "2023-10-02T00:00:00Z",
		Configuration: models.JobCreateRequest{
			Description: "Test job 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetJob(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJob() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeleteJob(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Delete existing job",
			args: args{
				id: "job1",
			},
		},
		{
			name: "Delete non-existing job", // Best Effort, doesn't error out
			args: args{
				id: "non-existing-job",
			},
		},
	}
	setup() // Call setup to clear the jobs map before each test
	// Add a test job to the jobs map
	jobs["job1"] = models.Job{
		ID:        "job1",
		Status:    string(JobStatusCreated),
		CreatedAt: "2023-10-01T00:00:00Z",
		Configuration: models.JobCreateRequest{
			Description: "Test job 1",
		},
	}
	// Add a second job to the jobs map
	jobs["job2"] = models.Job{
		ID:        "job2",
		Status:    string(JobStatusRunning),
		CreatedAt: "2023-10-02T00:00:00Z",
		Configuration: models.JobCreateRequest{
			Description: "Test job 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteJob(tt.args.id)
			_, exists := jobs[tt.args.id]
			if exists {
				t.Errorf("DeleteJob() did not delete job with id %v", tt.args.id)
			}
		})
	}
}
