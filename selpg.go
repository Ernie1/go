
/*=================================================================

Program name:
	selpg (SELect PaGes)

Purpose:
	Sometimes one needs to extract only a specified range of
pages from an input text file. This program allows the user to do
that.

===================================================================*/

package main

/*================================= imports =======================*/

import (
	"fmt"
	"os"
	"syscall"
	"os/exec"
	"io"
	"bufio"
	"strings"
	flag "github.com/ogier/pflag"
)

/*================================= types =========================*/

type sp_args struct {
	start_page  int
	end_page    int
	in_filename string
	page_len    int
	page_type   bool
	print_dest  string
}

/*================================= globals =======================*/

var sa = new(sp_args)

/*================================= init() ========================*/

func init() {
	flag.IntVarP(&sa.start_page, "s", "s", -1, "start page")
	flag.IntVarP(&sa.end_page, "e", "e", -1, "end page")
	flag.IntVarP(&sa.page_len, "l", "l", 72, "lines per page")
	flag.BoolVarP(&sa.page_type, "f", "f", false, "form-feed-delimited")
	flag.StringVarP(&sa.print_dest, "d", "d", "", "destination")
	flag.Usage = usage
}

/*================================= main()=========================*/

func main() {
	flag.Parse()
	process_args()
	process_input()
}

/*================================= process_args() ================*/

func process_args() {
	if flag.NFlag() < 1 {
	  /* handle mandatory args first */
	} else if sa.start_page < 1 {
		fmt.Fprintf(os.Stderr, "invalid start page %v\n", sa.start_page)
	} else if sa.end_page < 1 || sa.end_page < sa.start_page {
		fmt.Fprintf(os.Stderr, "invalid end page %v\n", sa.end_page)
	} else if !sa.page_type && sa.page_len < 1 {
		fmt.Fprintf(os.Stderr, "invalid page length %v\n", sa.page_len)
	} else {
		/* while there more args and they start with a '-' */
		if flag.NArg() > 0 {
			sa.in_filename = flag.Arg(0)
			if syscall.Access(sa.in_filename, syscall.O_RDONLY) != nil {
				fmt.Fprintf(os.Stderr, "input file \"%s\" does not exist or cannot be read\n", sa.in_filename)				
			} else { return }
		} else { return }
	}
	flag.Usage()
	os.Exit(1)
}

/*================================= process_input() ===============*/

func process_input() {
	fin := os.Stdin /* input stream */
	fout := os.Stdout /* output stream */
	line_ctr := 0 /* line counter */
	page_ctr := 1 /* page counter */
	var inpipe io.WriteCloser
	var err error

	/* set the input source */
	if sa.in_filename != "" {
		fin, err = os.Open(sa.in_filename)
	}

	/* set the output destination */
	if sa.print_dest != "" {
		cmd := exec.Command("lp", "-d", sa.print_dest)
		inpipe, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open pipe to \"%s\"\n", sa.print_dest)				
			flag.Usage()
			os.Exit(1)
		}
		cmd.Stdout = fout
		cmd.Start()
	}

	/* begin one of two main loops based on page type */
	if sa.page_type {
		reader := bufio.NewReader(fin)
		for {
			page, rerr := reader.ReadString('\f')
			if page_ctr >= sa.start_page {
				page = strings.Replace(page, "\f", "", -1)
				if sa.print_dest != "" {
					fmt.Fprintf(inpipe, "%s", page)
				} else {
					fmt.Fprintf(fout, "%s", page)
				}
			}
			page_ctr++
			if rerr == io.EOF || page_ctr > sa.end_page {
				break
			}
		}
	} else {
		line := bufio.NewScanner(fin)
		for line.Scan() {
			if page_ctr >= sa.start_page {
				if sa.print_dest != "" {
					fmt.Fprintf(inpipe, "%s\n", line.Text())
				} else {
					fmt.Fprintf(fout, "%s\n", line.Text())
				}
			}
			line_ctr++
			if line_ctr == sa.page_len {
				line_ctr = 0
				page_ctr++
				if page_ctr > sa.end_page {
					break
				}
			}
		}
	}

	/* end main loop */

	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr, "start_page (%d) greater than total pages (%d), no output written\n", sa.start_page, page_ctr)
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr, "end_page (%d) greater than total pages (%d), less output than expected\n", sa.end_page, page_ctr)
	} else { /* it was EOF, not error */
		fin.Close()
		if sa.print_dest != "" {
			inpipe.Close()
		}
		fmt.Fprintf(os.Stderr, "done\n");
	}
}

/*================================= usage() =======================*/

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: ./selpg [-s start_page] [-e end_page] [ -f | -l lines_per_page ] [ -d dest ] [ in_filename ]

Options:
`)
	flag.PrintDefaults()
}

/*================================= EOF ===========================*/
