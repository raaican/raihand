# Issues

## go.mod

```go
module github.com/raaican/raihand

go 1.26.2

require (
	charm.land/bubbles/v2 v2.1.0
	charm.land/bubbletea/v2 v2.0.6
	charm.land/lipgloss/v2 v2.0.3
	github.com/bwmarrin/discordgo v0.29.0
)
```

## go vet
```sh
 go vet .
# github.com/raaican/raihand/internal/ui
internal/ui/actions.go:256:21: cannot convert lipgloss.JoinVertical(lipgloss.Left, title, row, resultLine) (value of type string) to type tea.View
internal/ui/logs.go:25:24: cannot use 80 (untyped int constant) as viewport.Option value in argument to viewport.New
internal/ui/logs.go:25:28: cannot use 20 (untyped int constant) as viewport.Option value in argument to viewport.New
```

## staticcheck
```sh
 staticcheck .
-: # github.com/raaican/raihand/internal/ui
internal/ui/actions.go:256:21: cannot convert lipgloss.JoinVertical(lipgloss.Left, title, row, resultLine) (value of type string) to type tea.View
internal/ui/logs.go:25:24: cannot use 80 (untyped int constant) as viewport.Option value in argument to viewport.New
internal/ui/logs.go:25:28: cannot use 20 (untyped int constant) as viewport.Option value in argument to viewport.New (compile)
```

## file tree
```
.
├── go.mod
├── go.sum
├── internal
│   ├── bot
│   │   ├── bot.go
│   │   └── commands.go
│   └── ui
│       ├── actions.go
│       ├── dashboard.go
│       ├── guilds.go
│       ├── logs.go
│       ├── root.go
│       └── styles.go
├── issues.md
├── LICENSE
└── main.go

4 directories, 13 files
```
