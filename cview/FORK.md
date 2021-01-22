This document explains why [tview](https://github.com/rivo/tview) was forked to
create [cview](https://gitlab.com/tslocum/cview). It also explains any
differences between the projects and tracks which tview pull requests have been
merged into cview.

# Why fork?

[rivo](https://github.com/rivo), the creator and sole maintainer of tview,
explains his reviewing and merging process in a [GitHub comment](https://github.com/rivo/tview/pull/298#issuecomment-559373851).

He states that he does not have the necessary time or interest to review,
discuss and merge pull requests:

>this project is quite low in priority. It doesn't generate any income for me
>and, unfortunately, reviewing issues and PRs is also not much "fun".

>But some other people submitted large PRs which will cost me many hours to
>review. (I had to chuckle a bit when I saw [this comment](https://github.com/rivo/tview/pull/363#issuecomment-555484734).)

>Lastly, I'm the one who ends up maintaining this code. I have to be 100%
>behind it, understand it 100%, and be able to make changes to it later if
> necessary.

cview aims to solve these issues by increasing the number of project
maintainers and allowing code changes which may be outside of tview's scope.

# Differences

## cview is thread-safe

tview [is not thread-safe](https://godoc.org/github.com/rivo/tview#hdr-Concurrency).

## Application.QueueUpdate and Application.QueueUpdateDraw do not block

tview [blocks until the queued function returns](https://github.com/rivo/tview/blob/fe3052019536251fd145835dbaa225b33b7d3088/application.go#L510).

## Double clicks are not handled by default

All clicks are handled as single clicks until an interval is set with [Application.SetDoubleClickInterval](https://docs.rocketnine.space/gitlab.com/tslocum/cview/#Application.SetDoubleClickInterval).

## Lists and Forms do not wrap around by default

Call `SetWrapAround(true)` to wrap around when navigating.

# Merged pull requests

The following tview pull requests have been merged into cview:

- [#378 Throttle resize handling](https://github.com/rivo/tview/pull/378)
- [#368 Add support for displaying text next to a checkbox](https://github.com/rivo/tview/pull/368)
- [#363 Mouse support](https://github.com/rivo/tview/pull/363)
- [#353 Add window size change handler](https://github.com/rivo/tview/pull/353)
- [#347 Handle ANSI code 39 and 49](https://github.com/rivo/tview/pull/347)
- [#336 Don't skip regions at end of line](https://github.com/rivo/tview/pull/336)
- [#296 Fixed TextView's reset &#x5B;-&#x5D; setting the wrong color](https://github.com/rivo/tview/pull/296)
