# How to Contribute

CNI is [Apache 2.0 licensed](LICENSE) and accepts contributions via GitHub
pull requests. This document outlines some of the conventions on development
workflow, commit message formatting, contact points and other resources to make
it easier to get your contribution accepted.

We gratefully welcome improvements to documentation as well as to code.

# Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

## Contribution workflow

This is a rough outline of how to prepare a contribution:

- Create a topic branch from where you want to base your work
  (usually branched from main).
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- If you changed code:
   - Add unit and integration tests to cover your changes
     (see existing tests for style).
   - Run all the test and ensure they pass.
- Make sure any new code files have a license header.
- Submit a pull request to the original repository.
- Examine the check results on the PR. All required checks should pass.

## How to run the test suite
We generally require test coverage of any new features or bug fixes.

Unit test should not require any dependencies and may run on any platform.

Integration tests depend on nftables and require a Linux platform with nftable
support in the kernel and the `nft` executable.

To run the unit-tests, use this command from the project root:
```bash
./automation/run-tests.sh --unit-test
```

To run the integration-tests, use this command from the project root:
```bash
./automation/run-tests.sh --integration-test
```

## Help utilities

In order to help resolve the formatting checks, developers can use the
provided `./automation/go-format.sh` script to format inline the code.


# Acceptance policy

The following points will make a PR more likely to be accepted:

 - A well-described requirement.
 - Tests for new code.
 - Tests for old code!
 - A good commit message (see below).

In general, we will merge a PR once at least one maintainer has endorsed it.
For substantial changes, more people may become involved,
and you might get asked to resubmit the PR or divide the changes into more than
one PR/commit.

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

See https://chris.beams.io/posts/git-commit .
