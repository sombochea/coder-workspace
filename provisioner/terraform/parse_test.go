//go:build linux

package terraform

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/coder/coder/provisionersdk"
	"github.com/coder/coder/provisionersdk/proto"
	"github.com/stretchr/testify/require"
	"storj.io/drpc/drpcconn"
)

func TestParse(t *testing.T) {
	t.Parallel()

	// Create an in-memory provisioner to communicate with.
	client, server := provisionersdk.TransportPipe()
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
		cancelFunc()
	})
	go func() {
		err := Serve(ctx, &ServeOptions{
			ServeOptions: &provisionersdk.ServeOptions{
				Transport: server,
			},
		})
		require.NoError(t, err)
	}()
	api := proto.NewDRPCProvisionerClient(drpcconn.New(client))

	for _, tc := range []struct {
		Name     string
		Files    map[string]string
		Response *proto.Parse_Response
	}{{
		Name: "single-variable",
		Files: map[string]string{
			"main.tf": `variable "A" {
				description = "Testing!"
			}`,
		},
		Response: &proto.Parse_Response{
			ParameterSchemas: []*proto.ParameterSchema{{
				Name:        "A",
				Description: "Testing!",
			}},
		},
	}, {
		Name: "default-variable-value",
		Files: map[string]string{
			"main.tf": `variable "A" {
				default = "wow"
			}`,
		},
		Response: &proto.Parse_Response{
			ParameterSchemas: []*proto.ParameterSchema{{
				Name:         "A",
				DefaultValue: "\"wow\"",
			}},
		},
	}, {
		Name: "variable-validation",
		Files: map[string]string{
			"main.tf": `variable "A" {
				validation {
					condition = var.A == "value"
				}
			}`,
		},
		Response: &proto.Parse_Response{
			ParameterSchemas: []*proto.ParameterSchema{{
				Name:                "A",
				ValidationCondition: `var.A == "value"`,
			}},
		},
	}} {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			// Write all files to the temporary test directory.
			directory := t.TempDir()
			for path, content := range tc.Files {
				err := os.WriteFile(filepath.Join(directory, path), []byte(content), 0644)
				require.NoError(t, err)
			}

			response, err := api.Parse(ctx, &proto.Parse_Request{
				Directory: directory,
			})
			require.NoError(t, err)

			// Ensure the want and got are equivalent!
			want, err := json.Marshal(tc.Response)
			require.NoError(t, err)
			got, err := json.Marshal(response)
			require.NoError(t, err)

			require.Equal(t, string(want), string(got))
		})
	}
}