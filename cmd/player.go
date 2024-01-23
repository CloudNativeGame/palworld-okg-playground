package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

var kickPlayerCmd = &cobra.Command{
	Use:   "kick",
	Short: "Kick a player from a PalWorld gameserver",
	Long:  `Kick a player from a PalWorld gameserver`,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

func init() {
	kickPlayerCmd.Flags().StringP("name", "n", "", "Name of the player")
	playerCmd.AddCommand(kickPlayerCmd)
}
