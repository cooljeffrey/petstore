package main

import (
	"flag"
	"fmt"
	"github.com/cooljeffrey/petstore/model"
	"github.com/cooljeffrey/petstore/service"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
)

// Show usage info on command line
func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func main() {
	// command line args
	fs := flag.NewFlagSet("petstore", flag.ExitOnError)
	var (
		httpAddr = fs.String(
			"http-addr",
			"0.0.0.0",
			"HTTP address of petsore host")
		httpPort = fs.String(
			"http-port",
			"8080",
			"HTTP port")
		mongoUri = fs.String(
			"mongo-uri",
			"mongodb://localhost:27017",
			"mongo db uri")
		mongoDbName = fs.String(
			"mongo-dbname",
			"petstore",
			"mongo db name")
		mongoDbTimeoutSeconds = fs.Int64(
			"mongo-timeout", 10,
			"mongo db operatoin timeout in seconds")
		publicBaseUri = fs.String(
			"public-uri",
			"/images",
			"the base uri of public files")
		publicFilePath = fs.String(
			"public-path",
			"./public",
			"the folder basing on current working dir")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fs.Usage()
		os.Exit(1)
	}

	// init logger
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stderr)

	// init storage
	storage, err := model.NewMongoStorage(*mongoUri, *mongoDbName, *mongoDbTimeoutSeconds, logger)
	if err != nil {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	// init services
	services := service.Services{
		UserService: service.NewUserService(log.WithPrefix(logger, "service", "user"), storage),
		PetService: service.NewPetService(
			log.WithPrefix(logger, "service", "pet"), storage, *publicBaseUri, *publicFilePath),
		StoreService: service.NewStoreService(log.WithPrefix(logger, "service", "store"), storage),
	}

	// init routes
	r := service.SetupRoutes(&services, log.WithPrefix(logger, "service", "routing"))

	// format server address
	addr := fmt.Sprintf("%s:%s", *httpAddr, *httpPort)

	// catch http server error
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// catch transport error
	go func() {
		_ = logger.Log("transport", "HTTP", "addr", addr)
		errs <- http.ListenAndServe(addr, r)
	}()

	// output (error) event on exit
	_ = logger.Log("exit", <-errs)
}
