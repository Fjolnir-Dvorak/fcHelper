// Copyright © 2017 Raphael Tiersch <tiersch.raphael@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package generate

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var (
	defaultMan = ""
)

// manCmd represents the man command
var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Generates a manual for this program",
	Long:  `Man which can be used as manpage.`,
	Run:   doMan,
}

func doMan(cmd *cobra.Command, args []string) {
	os.Mkdir(defaultMan, os.ModePerm)
	doc.GenManTree(cmd.Root(), &doc.GenManHeader{Title: "fcHelper", Section: "man"}, defaultMan)
}

func init() {
	GenerateCmd.AddCommand(manCmd)
	manCmd.Flags().StringVarP(&defaultMan, "destination", "d",
		defaultMan, "Output location.")
}
