//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cmd

import (
	"fmt"
	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/view"
	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(serviceInspectCmd)
}

const serviceInspectExample = `
  # Get information for 'redis' service in 'ns-demo' namespace
  lb service inspect ns-demo redis
`

var serviceInspectCmd = &cobra.Command{
	Use:     "inspect [NAMESPACE]/[NAME]",
	Short:   "Service info by name",
	Example: serviceInspectExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace, name, err := serviceParseSelfLink(args[0])
		checkError(err)

		cli := envs.Get().GetClient()
		svc, err := cli.Cluster.V1().Namespace(namespace).Service(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		routes, err := cli.Cluster.V1().Namespace(namespace).Route().List(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, r := range *routes {
			for _, rule := range r.Spec.Rules {
				if rule.Service == svc.Meta.Name {
					fmt.Println("exposed:", r.Status.State, r.Spec.Domain, r.Spec.Port)
				}
			}

		}

		ss := view.FromApiServiceView(svc)
		ss.Print()
	},
}
