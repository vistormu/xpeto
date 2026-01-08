<a name="readme-top"></a>

<div align="center">

<a href="https://github.com/vistormu/xpeto" target="_blank" title="go to the repo"><img width="196px" alt="xpeto logo" src="/assets/xpeto.png"></a>

# xpeto<br>an ECS runtime for apps and games!

_xpeto_ is a small ECS based runtime written in Go, designed for real time simulation and simple 2D applications.

<br>

[![go version][go_version_img]][go_dev_url]
[![License][repo_license_img]][repo_license_url]

<br>

</div>

> [!WARNING]
> this is a personal project under active development
> expect breaking changes, slow development and bugs

## overview

_xpeto_ focuses on **execution clarity** rather than tooling, editors, or complex asset pipelines. you describe *what exists* and *when systems run*. the engine takes care of running it.

_xpeto_ works both with and without rendering. this makes it suitable for interactive apps, simulations, and headless workloads.

---

## what xpeto is good at

- deterministic update and fixed step simulation
- clear system ordering through explicit stages
- headless execution for servers, tests, and experiments
- plugin based architecture with minimal coupling
- simple 2D rendering provided by a backend

typical use cases include:
- simulation driven games
- robotics and control experiments
- digital twins and interactive visualisations
- research and teaching projects

---

## what xpeto is not

xpeto is not:
- a general purpose game engine
- an editor driven framework
- a full physics, networking, or audio solution
- a drop in replacement for Unity, Godot, or Unreal

those features can exist as plugins, but they are not the goal of the core

---

## core ideas

- **ECS first**  
  entities are identifiers. components are plain data. systems are scheduled functions.

- **explicit execution**  
  systems run in stages. ordering is visible and configurable.

- **plugins over globals**  
  features are added through plugins that register systems and resources.

- **headless friendly**  
  rendering and input live in backends, not in the core.

- **simple Go code**  
  no code generation, no macros, no hidden control flow.

---

## architecture in a nutshell

xpeto is organised into four layers:

1. **core scaffold**  
   ECS, world, scheduler, and resources.

2. **core plugins**  
   logging, events, clock, and window intent.

3. **default plugins**  
   assets, rendering, sprites, shapes, text, input.

4. **backends**  
   ebiten for interactive apps and a headless backend for non visual execution.

the backend provides the main loop and platform integration. everything else is portable.

---

## status

_xpeto_ is an experimental project developed alongside research work.

the core concepts are stable, but APIs may change as new use cases appear.

---

## license

MIT

## should I use xpeto?

| if you want to…                                  | _xpeto_ |
|--------------------------------------------------|:-------:|
| control exactly when systems run                 | yes     |
| run the same logic with or without rendering     | yes     |
| write everything in plain Go                     | yes     |
| avoid hidden engine control flow                 | yes     |
| build small to medium 2D apps or games           | yes     |
| run simulations or experiments headless          | yes     |
| use an editor or visual scene tools              | no      |
| get a full physics or networking stack           | no      |
| rely on a large asset and plugin ecosystem       | no      |

## examples

check this snippet for getting the _xpeto_ vibes:

```go
package main

import (
    "fmt"

    "github.com/vistormu/xpeto"
    "github.com/vistormu/xpeto/backends/headless"
)

func helloAndExit(w *xp.World) {
    fmt.Println("hello, world!")
    xp.AddEvent(w, xp.ExitAppEvent{})
}

func Pkg(w *xp.World, sch *xp.Scheduler) {
   xp.AddSystem(sch, xp.Update, helloAndExit)
}

func main() {
    xp.NewApp(headless.Backend,
        xp.AppOpt.Pkgs(
            Pkg,
        )
    ).Run()
}
```

more examples are avialable [here](/examples)


[go_version_img]: https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go
[go_dev_url]: https://go.dev/
[go_report_img]: https://goreportcard.com/badge/github.com/vistormu/xpeto
[go_report_url]: https://goreportcard.com/report/github.com/vistormu/xpeto
[repo_license_img]: https://img.shields.io/github/license/vistormu/xpeto?style=for-the-badge
[repo_license_url]: https://github.com/vistormu/xpeto/blob/main/LICENSE
