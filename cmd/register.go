/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"

	"github.com/erik-farmer/me-and-u/data"
	"github.com/spf13/cobra"
)

var username string
var password string

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Adds a user to the database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register called")
		db, err := data.NewDB()
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()

		err = data.CreateUser(db, username, password)
		if err != nil {
			log.Fatalf("Failed to register user: %v", err)
		}
		fmt.Printf("User '%s' registered successfully\n", username)
	},
}

func init() {
	cobra.OnInitialize(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	})
	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
