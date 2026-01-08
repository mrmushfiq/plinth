package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "cluster":
		handleClusterCommand()
	case "nodes":
		handleNodesCommand()
	case "repair":
		handleRepairCommand()
	case "object":
		handleObjectCommand()
	case "costs":
		handleCostsCommand()
	case "version":
		fmt.Println("objctl version 0.1.0-alpha")
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`objctl - Plinth Admin CLI

Usage:
  objctl <command> [arguments]

Commands:
  cluster    Cluster management
    status     Show cluster status
  
  nodes      Node management
    list       List all nodes
    status     Show node status
  
  repair     Repair operations
    status     Show repair worker status
    run        Trigger repair cycle
  
  object     Object operations
    stat       Show object metadata
    verify     Verify object integrity
  
  costs      Cost tracking
    bucket     Show costs for a bucket
    top        Show top objects by cost
  
  version    Show version
  help       Show this help message

Examples:
  objctl cluster status
  objctl nodes list
  objctl repair status
  objctl object stat bucket/key
  objctl costs bucket ml-datasets

For more information, visit: https://github.com/mrmushfiq/plinth`)
}

func handleClusterCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: objctl cluster <status>")
		return
	}
	subcmd := os.Args[2]
	if subcmd == "status" {
		// TODO: Implement cluster status
		fmt.Println("Cluster status: healthy (TODO: implement)")
	}
}

func handleNodesCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: objctl nodes <list|status>")
		return
	}
	subcmd := os.Args[2]
	switch subcmd {
	case "list":
		// TODO: Implement node list
		fmt.Println("Listing nodes... (TODO: implement)")
	case "status":
		// TODO: Implement node status
		fmt.Println("Node status... (TODO: implement)")
	}
}

func handleRepairCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: objctl repair <status|run>")
		return
	}
	subcmd := os.Args[2]
	switch subcmd {
	case "status":
		// TODO: Implement repair status
		fmt.Println("Repair status... (TODO: implement)")
	case "run":
		// TODO: Trigger repair cycle
		fmt.Println("Triggering repair cycle... (TODO: implement)")
	}
}

func handleObjectCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: objctl object <stat|verify> <bucket/key>")
		return
	}
	subcmd := os.Args[2]
	switch subcmd {
	case "stat":
		if len(os.Args) < 4 {
			fmt.Println("Usage: objctl object stat <bucket/key>")
			return
		}
		// TODO: Show object metadata
		fmt.Printf("Object stats for %s... (TODO: implement)\n", os.Args[3])
	case "verify":
		if len(os.Args) < 4 {
			fmt.Println("Usage: objctl object verify <bucket/key>")
			return
		}
		// TODO: Verify object
		fmt.Printf("Verifying %s... (TODO: implement)\n", os.Args[3])
	}
}

func handleCostsCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: objctl costs <bucket|top> [args]")
		return
	}
	subcmd := os.Args[2]
	switch subcmd {
	case "bucket":
		if len(os.Args) < 4 {
			fmt.Println("Usage: objctl costs bucket <bucket-name>")
			return
		}
		// TODO: Show bucket costs
		fmt.Printf("Costs for bucket %s... (TODO: implement)\n", os.Args[3])
	case "top":
		// TODO: Show top objects by cost
		fmt.Println("Top objects by cost... (TODO: implement)")
	}
}
