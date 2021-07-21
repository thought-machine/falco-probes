# Contributing Guidance
## Appreciated Contributions
There are currently two ways that you can contribute:
* [ Issues](https://github.com/thought-machine/falco-probes/issues/new/choose) - for discussion of suspected bugs, suggested improvements or feature requests. 
* [Pull Requests](https://github.com/thought-machine/falco-probes/pulls) - for fixes and improvements.

## How to Contribute
Before contributing, please familiarise yourself with our [Code of Conduct](https://github.com/thought-machine/falco-probes/blob/master/CODE_OF_CONDUCT.md) & the terms of our [License](https://github.com/thought-machine/falco-probes/blob/master/LICENSE).
### Issues
1. Create a [new issue](https://github.com/thought-machine/falco-probes/issues/new/choose), please apply at least 1 label.
2. Once discussed and resolved a Thought Machine Organisation member will close the issue. 
### Pull Requests
Pull requests to this project should either:
- Link to a new / existing Issue.
- Be initiated from a ticket from the Thought Machine jira board, the link of the PR should be commented on the jira ticket.

The process is as follows:
1. Clone the repository using SSH `git clone git@github.com:thought-machine/falco-probes.git`.
2. Create a local feature branch `git checkout -b <feature-branch`, and set it's upstream `git branch -u origin`.
3. Add your commits with a [descriptive](https://chris.beams.io/posts/git-commit/) subject line and body using `git commit -a`.
4. Push your commits to a remote feature branch `git push origin <feature-branch>`.
5. Create a pull request through the UI, please request a review from the Thought Machine Organisation members and apply any relevant labels.
6. Once discussed and approved, the UI will be used to squash and merge the PR onto master by a Thought Machine Organisation member.
7. Once merged the remote branch should be deleted.
## Future Considerations
- Should interest in contributing increase, the Thought Machine Organisation members may consider implementing a CLA/DCO as needed.

## Working in the falco-probes repository
Here's some information that might be helpful while working on PRs:
- The repository layout largely maps to the structure of [golang-standards/project-layout](https://github.com/golang-standards/project-layout#standard-go-project-layout).
- The output of our [Releases](https://github.com/thought-machine/falco-probes/releases) are the Falco probe files published as assets (as detailed in [REPOSITORY_DESIGN.md](https://github.com/thought-machine/falco-probes/blob/master/docs/REPOSITORY_DESIGN.md)), rather a package of this repository's code.
- The automated [workflows](https://github.com/thought-machine/falco-probes/tree/master/.github/workflows) involved to create each [Release](https://github.com/thought-machine/falco-probes/releases) can be seen in [Actions](https://github.com/thought-machine/falco-probes/actions).