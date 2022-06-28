# Contributing Guidelines

Thank you for your interest in contributing to apkstat! Here are a few resources
to help you get started.

For now, the [#accrescent:matrix.org] Matrix room is the best place to get help
with development. [@lberrymage:matrix.org] is the lead developer and you are
welcome to DM him as well.

Most documentation is pulled from the respecitve code in AOSP. You can find said
code at
https://android.googlesource.com/platform/frameworks/base/+/master/libs/androidfw.
The most relevant files are `ResourceTypes.cpp` and `include/ResourceTypes.h`,
although you may need to reference others at times.

If you're translating AOSP code into apkstat, keep the code as similar as
possible to upstream (i.e. don't extensively reorder/optimize it). This makes it
much easier to check what work still needs to be done on apkstat and to adapt to
upstream changes.

When adding fields to the Android manifest, please verify they can be parsed
correctly by testing on multiple APKs.

The project's coding style and conventions are outlined below. Please check your
branch against them before making a PR to expeditide the review process.

## Code style

- Wrap lines at 100 columns. This isn't a hard limit, but will be enforced
  unless wrapping a line looks uglier than extending it by a few columns.
- Format with `gofmt -w -s .`
- Where `gofmt` doesn't have a preference for style, try to be consistent with
  the rest of the codebase.

## Code conventions

- Do not use any unsafe code. This includes the "C" and "unsafe" packages.
- Avoid third-party libraries. `apkstat` currently only uses the standard
  library and has little reason to use anything else, so a compelling argument
  would need to be made for it to use anything outside of the standard library.
- When translating AOSP code, keep it as similar as possible to upstream to make
  referencing and review easier.

## Vulnerability reports

Report all vulnerabilities in accordance with Accrescent's [security policy].

[#accrescent:matrix.org]: https://matrix.to/#/#accrescent:matrix.org
[@lberrymage:matrix.org]: https://matrix.to/#/@lberrymage:matrix.org
[security policy]: SECURITY.md
