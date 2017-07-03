// Copyright © 2017 Walter Scheper <walter.scheper@gmal.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wfscheper/fated/fate"
)

// drawCmd represents the draw command
var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Draw a fate card",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if foreground {
			fmt.Println("Press enter to draw")
			for {
				reader := bufio.NewReader(os.Stdin)
				reader.ReadString('\n')
				rolls := fate.RollDice(4)
				fmt.Println(fate.RenderCard(rolls))
			}
		} else {
			rolls := fate.RollDice(4)
			fmt.Println(fate.RenderCard(rolls))
		}
	},
}

func init() {
	RootCmd.AddCommand(drawCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// drawCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// drawCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
