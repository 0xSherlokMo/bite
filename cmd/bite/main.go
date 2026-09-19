// Command bite automates delivery apps over their private mobile APIs.
//
// Usage:
//
//	bite talabat session import <headers.json>
//	bite talabat profile
//	bite talabat addresses
//	bite talabat food restaurants --lat <f> --lon <f>
//	bite talabat food menu <branchId>
//	bite talabat mart cart --vendor <uuid> --branch <id> --chain <id> --lat <f> --lon <f>
//	bite talabat mart add  --vendor <uuid> --branch <id> --chain <id> --product <uuid> --qty <n> --lat <f> --lon <f>
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/salah/bite/internal/providers/talabat"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "talabat":
		runTalabat(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown provider %q\n", os.Args[1])
		usage()
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `bite — automate delivery apps

  bite talabat session import <headers.json>
  bite talabat profile
  bite talabat addresses
  bite talabat food restaurants --lat <f> --lon <f>
  bite talabat food menu <branchId>
  bite talabat mart cart --vendor <uuid> --branch <id> --chain <id> --lat <f> --lon <f>
  bite talabat mart add  --vendor <uuid> --branch <id> --chain <id> --product <uuid> [--qty 1] --lat <f> --lon <f>
`)
	os.Exit(2)
}

func runTalabat(args []string) {
	if len(args) == 0 {
		usage()
	}
	ctx := context.Background()

	// session import needs no client.
	if args[0] == "session" {
		if len(args) < 3 || args[1] != "import" {
			die("usage: bite talabat session import <headers.json>")
		}
		loc, err := talabat.ImportHeaders(args[2])
		check(err)
		fmt.Printf("session imported -> %s\n", loc)
		return
	}

	client, err := talabat.New()
	check(err)

	switch args[0] {
	case "profile":
		p, err := client.Profile(ctx)
		check(err)
		printJSON(p)

	case "addresses":
		a, err := client.Addresses(ctx)
		check(err)
		printJSON(a)

	case "food":
		runFood(ctx, client, args[1:])

	case "mart":
		runMart(ctx, client, args[1:])

	default:
		die("unknown talabat command: " + args[0])
	}
}

func runFood(ctx context.Context, c *talabat.Client, args []string) {
	if len(args) == 0 {
		die("usage: bite talabat food [restaurants|menu]")
	}
	switch args[0] {
	case "restaurants":
		fs := flag.NewFlagSet("restaurants", flag.ExitOnError)
		lat := fs.Float64("lat", 0, "latitude")
		lon := fs.Float64("lon", 0, "longitude")
		area := fs.Int("area", 0, "area id (from `bite talabat addresses`)")
		page := fs.Int("page", 1, "page (starts at 1)")
		_ = fs.Parse(args[1:])
		rs, err := c.Restaurants(ctx, *lat, *lon, *area, *page)
		check(err)
		for _, r := range rs {
			fmt.Printf("%-8d branch=%-8d  %-32s  ★%.1f (%d)  %s\n",
				r.VendorID, r.BranchID, trunc(r.Name, 32), r.Rating, r.RatingsCount, r.DeliveryText)
		}
		fmt.Printf("\n%d restaurants\n", len(rs))

	case "menu":
		if len(args) < 2 {
			die("usage: bite talabat food menu <branchId>")
		}
		var branch int
		if _, err := fmt.Sscanf(args[1], "%d", &branch); err != nil {
			die("branchId must be numeric")
		}
		m, err := c.Menu(ctx, branch)
		check(err)
		printJSON(m)

	default:
		die("unknown food command: " + args[0])
	}
}

func runMart(ctx context.Context, c *talabat.Client, args []string) {
	if len(args) == 0 {
		die("usage: bite talabat mart [cart|add]")
	}
	fs := flag.NewFlagSet("mart", flag.ExitOnError)
	vendor := fs.String("vendor", "", "dh_vendor_id (uuid)")
	branch := fs.String("branch", "", "branch_id")
	chain := fs.String("chain", "", "chain_id")
	lat := fs.Float64("lat", 0, "latitude")
	lon := fs.Float64("lon", 0, "longitude")
	product := fs.String("product", "", "product_id (uuid) for add")
	qty := fs.Int("qty", 1, "quantity for add")
	_ = fs.Parse(args[1:])

	ref := talabat.VendorRef{DHVendorID: *vendor, BranchID: *branch, ChainID: *chain, Lat: *lat, Lon: *lon}

	switch args[0] {
	case "cart":
		cart, err := c.GetCart(ctx, ref)
		check(err)
		printJSON(cart)
	case "add":
		if *product == "" {
			die("mart add requires --product")
		}
		cart, err := c.AddToCart(ctx, ref, *product, *qty)
		check(err)
		printJSON(cart)
	default:
		die("unknown mart command: " + args[0])
	}
}

func printJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, "error: "+msg)
	os.Exit(2)
}

func check(err error) {
	if err == nil {
		return
	}
	if ae, ok := err.(*talabat.APIError); ok && ae.Expired() {
		fmt.Fprintln(os.Stderr, "error: session expired — re-import a fresh capture with `bite talabat session import <file>`")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
