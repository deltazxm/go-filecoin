package commands

import (
	"fmt"
	"io"
	"time"

	cmds "gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	peer "gx/ipfs/QmQsErDt8Qgw1XrsXf2BpEzDgGWtB1YLsTAARBup5b6B9W/go-libp2p-peer"
	cmdkit "gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"

	"github.com/filecoin-project/go-filecoin/api"
)

var pingCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Send echo request packets to p2p network members",
		ShortDescription: `
'ping' is a tool to test sending data to other nodes. It finds nodes
via the routing system, sends pings, waits for pongs, and prints out round-
trip latency information.
		`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("peer ID", true, true, "ID of peer to be pinged.").EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmdkit.UintOption("count", "n", "Number of ping messages to send.").WithDefault(10),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) {
		peerID, err := peer.IDB58Decode(req.Arguments[0])
		if err != nil {
			err = fmt.Errorf("failed to parse peer address '%s': %s", req.Arguments[0], err)
			re.SetError(err, cmdkit.ErrNormal)
			return
		}

		numPings, _ := req.Options["count"].(uint)

		ch, err := GetAPI(env).Ping().Ping(req.Context, peerID, numPings, time.Second)
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal)
			return
		}

		for p := range ch {
			re.Emit(p) // nolint: errcheck
		}
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, p *api.PingResult) error {
			if len(p.Text) > 0 {
				fmt.Fprintln(w, p.Text) // nolint: errcheck
			} else if p.Success {
				fmt.Fprintf(w, "Pong received: time=%.2f ms\n", p.Time.Seconds()*1000) // nolint: errcheck
			} else {
				fmt.Fprintf(w, "Pong failed\n") // nolint: errcheck
			}
			return nil
		}),
	},
	Type: api.PingResult{},
}
