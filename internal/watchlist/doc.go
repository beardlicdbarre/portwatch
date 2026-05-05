// Package watchlist provides a priority watch set for ports that are
// expected to remain open at all times.
//
// # Overview
//
// A Watchlist is created from a slice of Entry values, each describing a
// protocol and port number (e.g. tcp:443). During every scan cycle the
// daemon calls Missing to find any watched ports that have disappeared
// from the current snapshot. Such disappearances are surfaced as critical
// alerts independently of the normal filter and rate-limit pipeline.
//
// # Usage
//
//	wl := watchlist.New([]watchlist.Entry{
//		{Proto: "tcp", Port: 443},
//		{Proto: "tcp", Port: 22},
//	})
//
//	if gone := wl.Missing(currentPorts); len(gone) > 0 {
//		// emit high-priority alert for each entry in gone
//	}
//
// Protocol matching is case-insensitive; "TCP" and "tcp" are equivalent.
package watchlist
