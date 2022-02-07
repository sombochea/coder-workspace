package coderd_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/coder/coder/coderd"
	"github.com/coder/coder/coderd/coderdtest"
	"github.com/coder/coder/database"
	"github.com/coder/coder/provisioner/echo"
	"github.com/coder/coder/provisionersdk/proto"
)

func TestWorkspaceHistoryLogsByName(t *testing.T) {
	t.Parallel()
	t.Run("List", func(t *testing.T) {
		t.Parallel()
		client := coderdtest.New(t)
		user := coderdtest.CreateInitialUser(t, client)
		coderdtest.NewProvisionerDaemon(t, client)
		project := coderdtest.CreateProject(t, client, user.Organization)
		version := coderdtest.CreateProjectVersion(t, client, user.Organization, project.Name, &echo.Responses{
			Parse: echo.ParseComplete,
			Provision: []*proto.Provision_Response{{
				Type: &proto.Provision_Response_Log{
					Log: &proto.Log{
						Level:  proto.LogLevel_INFO,
						Output: "log-output",
					},
				},
			}, {
				Type: &proto.Provision_Response_Complete{
					Complete: &proto.Provision_Complete{},
				},
			}},
		})
		coderdtest.AwaitProjectVersionImported(t, client, user.Organization, project.Name, version.Name)
		workspace := coderdtest.CreateWorkspace(t, client, "me", project.ID)
		history, err := client.CreateWorkspaceHistory(context.Background(), "", workspace.Name, coderd.CreateWorkspaceHistoryRequest{
			ProjectVersionID: version.ID,
			Transition:       database.WorkspaceTransitionCreate,
		})
		require.NoError(t, err)

		// Successfully return empty logs before the job starts!
		logs, err := client.WorkspaceHistoryLogs(context.Background(), "", workspace.Name, history.Name)
		require.NoError(t, err)
		require.NotNil(t, logs)
		require.Len(t, logs, 0)

		coderdtest.AwaitWorkspaceHistoryProvisioned(t, client, "", workspace.Name, history.Name)

		// Return the log after completion!
		logs, err = client.WorkspaceHistoryLogs(context.Background(), "", workspace.Name, history.Name)
		require.NoError(t, err)
		require.NotNil(t, logs)
		require.Len(t, logs, 1)
	})

	t.Run("StreamAfterComplete", func(t *testing.T) {
		t.Parallel()
		client := coderdtest.New(t)
		user := coderdtest.CreateInitialUser(t, client)
		coderdtest.NewProvisionerDaemon(t, client)
		project := coderdtest.CreateProject(t, client, user.Organization)
		version := coderdtest.CreateProjectVersion(t, client, user.Organization, project.Name, &echo.Responses{
			Parse: echo.ParseComplete,
			Provision: []*proto.Provision_Response{{
				Type: &proto.Provision_Response_Log{
					Log: &proto.Log{
						Level:  proto.LogLevel_INFO,
						Output: "log-output",
					},
				},
			}, {
				Type: &proto.Provision_Response_Complete{
					Complete: &proto.Provision_Complete{},
				},
			}},
		})
		coderdtest.AwaitProjectVersionImported(t, client, user.Organization, project.Name, version.Name)
		workspace := coderdtest.CreateWorkspace(t, client, "me", project.ID)
		before := time.Now().UTC()
		history, err := client.CreateWorkspaceHistory(context.Background(), "", workspace.Name, coderd.CreateWorkspaceHistoryRequest{
			ProjectVersionID: version.ID,
			Transition:       database.WorkspaceTransitionCreate,
		})
		require.NoError(t, err)
		coderdtest.AwaitWorkspaceHistoryProvisioned(t, client, "", workspace.Name, history.Name)

		logs, err := client.FollowWorkspaceHistoryLogsAfter(context.Background(), "", workspace.Name, history.Name, before)
		require.NoError(t, err)
		log := <-logs
		require.Equal(t, "log-output", log.Output)
		// Make sure the channel automatically closes!
		_, ok := <-logs
		require.False(t, ok)
	})

	t.Run("StreamWhileRunning", func(t *testing.T) {
		t.Parallel()
		client := coderdtest.New(t)
		user := coderdtest.CreateInitialUser(t, client)
		coderdtest.NewProvisionerDaemon(t, client)
		project := coderdtest.CreateProject(t, client, user.Organization)
		version := coderdtest.CreateProjectVersion(t, client, user.Organization, project.Name, &echo.Responses{
			Parse: echo.ParseComplete,
			Provision: []*proto.Provision_Response{{
				Type: &proto.Provision_Response_Log{
					Log: &proto.Log{
						Level:  proto.LogLevel_INFO,
						Output: "log-output",
					},
				},
			}, {
				Type: &proto.Provision_Response_Complete{
					Complete: &proto.Provision_Complete{},
				},
			}},
		})
		coderdtest.AwaitProjectVersionImported(t, client, user.Organization, project.Name, version.Name)
		workspace := coderdtest.CreateWorkspace(t, client, "me", project.ID)
		history, err := client.CreateWorkspaceHistory(context.Background(), "", workspace.Name, coderd.CreateWorkspaceHistoryRequest{
			ProjectVersionID: version.ID,
			Transition:       database.WorkspaceTransitionCreate,
		})
		require.NoError(t, err)

		logs, err := client.FollowWorkspaceHistoryLogsAfter(context.Background(), "", workspace.Name, history.Name, time.Time{})
		require.NoError(t, err)
		log := <-logs
		require.Equal(t, "log-output", log.Output)
		// Make sure the channel automatically closes!
		_, ok := <-logs
		require.False(t, ok)
	})
}