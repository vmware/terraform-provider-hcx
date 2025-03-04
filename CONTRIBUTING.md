# Contributing to terraform-provider-hcx

Before you start working with terraform-provider-hcx please read our Developer Certificate
of Origin. All contributions to this repository must be signed as described on
that page. Your signature certifies that you wrote the patch or have the right
to pass it on as an open-source patch

For any questions about the DCO process, please refer to our [FAQ](https://cla.vmware.com/dco).

## Contribution Flow

This is a general outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work
- Make commits of logical units
- Make sure your commit messages are in the proper format (see below)
- Push your changes to a topic branch in your fork of the repository
- Submit a pull request

Example:

``` shell
git remote add upstream https://github.com/vmware/terraform-provider-hcx.git
git checkout -b my-new-feature main
git commit -a
git push origin my-new-feature
```

### Staying In Sync With Upstream

If your branch gets out of sync with the `vmware/main` branch, use the
following to update:

``` shell
git checkout my-new-feature
git fetch -a
git pull --rebase upstream main
git push --force-with-lease origin my-new-feature
```

### Updating Pull Requests

If your pull request fails to pass CI or needs changes based on code review,
you'll most likely want to squash these changes into existing commits.

If your pull request contains a single commit or your changes are related to
the most recent commit, you can simply amend the commit.

``` shell
git add .
git commit --amend
git push --force-with-lease origin my-new-feature
```

If you need to squash changes into an earlier commit, you can use:

``` shell
git add .
git commit --fixup <commit>
git rebase -i --autosquash main
git push --force-with-lease origin my-new-feature
```

Be sure to add a comment to the pull request indicating your new changes are
ready to review. GitHub does not provide a notification for a git push.

### Code Style

### Formatting Commit Messages

We follow the conventions on [How to Write a Git Commit Message](http://chris.beams.io/posts/git-commit/).

Be sure to include any related GitHub issue references in the commit message.
See [GFM syntax](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown)
for referencing issues and commits.

## Reporting Bugs and Creating Issues

When opening a new issue, try to roughly follow the commit message format
conventions above.