package main

import "github.com/happyxcj/wlog"

func main() {
	wlog.UseGlobalLogger(nil)
	defer wlog.Close()

	wlog.With(wlog.String("username", "root"),
		wlog.String("password", "root")).
		Info("failed to login account")

	wlog.Infow("failed to login account",
		wlog.String("username", "root"),
		wlog.String("password", "root"))

}
