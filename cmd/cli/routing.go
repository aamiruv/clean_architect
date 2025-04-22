package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	v2 "github.com/amirzayi/clean_architect/api/http/handler/v2"
	"github.com/amirzayi/rahjoo"
	"github.com/spf13/cobra"
)

var routingCmd = &cobra.Command{
	Use:   "routes",
	Short: "web application's route list",
	Run: func(cmd *cobra.Command, args []string) {
		routeList()
	},
}

const projectName = "github.com/amirzayi/clean_architect"

func routeList() {
	userV2Routes := v2.UserRoutes(nil, nil, nil)
	authV2Routes := v2.AuthRoutes(nil)

	routes := rahjoo.MergeRoutes(userV2Routes, authV2Routes)

	fmt.Println("--------------------------------------------------")
	fmt.Println("|  Route  |  Method  |  Handler  |  Middlewares  |")
	fmt.Println("--------------------------------------------------")

	getFuncName := func(i any) string {
		v := reflect.ValueOf(i)
		funcName := runtime.FuncForPC(v.Pointer()).Name()
		cleaned := strings.TrimPrefix(funcName, projectName)
		cleaned = strings.TrimSuffix(cleaned, "-fm")
		cleaned = strings.TrimSuffix(cleaned, ".func1")
		return cleaned
	}

	for path, methods := range routes {
		for method, action := range methods {
			fmt.Printf("| %s | %s ", path, method)
			fmt.Printf("| %s |", getFuncName(action.Handler()))
			middlewares := action.Middlewares()
			for _, middleware := range middlewares {
				fmt.Printf(" %s ", getFuncName(middleware))
			}
			fmt.Print("|\n")
		}
		fmt.Println()
	}
}
