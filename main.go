/*
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/mdhender/fhdb/config"
	"github.com/mdhender/fhdb/handlers"
	"github.com/mdhender/fhdb/store/jsondb"
	"github.com/mdhender/fhdb/way"
	"log"
	"mime"
	"net"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // force logs to be UTC

	// go depends on the operating system to associate extensions with mime-types.
	// the default works mostly for CSS, but this forces it.
	if err := mime.AddExtensionType(".css", "text/css; charset=utf-8"); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	cfg := config.Default()
	err := cfg.Load()
	if err == nil {
		err = run(cfg)
	}
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}

func run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("missing configuration information")
	}
	if cfg.Server.JWT.Key == "" || len(cfg.Server.JWT.Key) < 16 {
		return fmt.Errorf("jwt key length should be at least 16")
	}

	s := &Server{
		Router: way.NewRouter(),
	}
	s.Addr = net.JoinHostPort(cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port))
	s.IdleTimeout = cfg.Server.Timeout.Idle
	s.ReadTimeout = cfg.Server.Timeout.Read
	s.WriteTimeout = cfg.Server.Timeout.Write
	s.MaxHeaderBytes = 1 << 20 // TODO: make this configurable
	s.Data = cfg.Data
	var err error
	s.jdb, err = jsondb.Read(filepath.Join(cfg.Data, "galaxy.json"))
	if err != nil {
		return err
	}
	s.Handler = handlers.CORS(handlers.Version(s.Router, s.jdb.Version))
	//err = s.jdb.Write(filepath.Join(cfg.Data, "cluster.json"))
	//if err != nil {
	//	return err
	//}

	if err := s.Routes(cfg); err != nil {
		return err
	}

	if cfg.Server.TLS.Serve {
		log.Printf("[main] serving TLS on %s\n", s.Addr)
		return s.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
	}
	log.Printf("[main] listening on %s\n", s.Addr)
	return s.ListenAndServe()
}
