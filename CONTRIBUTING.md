### Which issue should I choose?

Currently any issue labelled `merged to devel` is currently in the pending release.
Any other issue is good-to-go, but if you're wanting to contribute, I'd advise choosing 
something that seems interesting. Unless something has been specifically claimed in
the comments section of the issue, it's fair game. Moreover, targets which can be 
prioritized will often be under the "[WIP] Release *.*.*" Pull request. An example:
https://github.com/GrappigPanda/notorious/pull/113


Don't ever be afraid to ask for help or clarification with an issue. I have the bad
habit of not always embellishing on particular issues, so you might find a cryptic
issue title with no further information. I'm trying to improve about that, though!


### How PRs are done for Notorious

Whenever an issue is being tackled for Notorious, the programmer ought to
checkout the origin/devel[1] branch into a local branch named `issue-<issue_number>`.


Once the changes have been made, unit tests have been added for newly adde features,
and any other housekeeping associated with finishing an issue have been completed,
the programmer should push to their remote repository.


Once your local repository changes are on your forked branch you can create a pull request.
It's important to note that all pull requests for Notorious should target the `devel` branch.
Any pull requests not targetting devel will be closed and I will ask to have the PR 
reopened for the `devel` branch. I genuinely wish that Github would allow changing 
targetted branch (hey, gitlab does it, so y'all should probably get on that), but it's
currently not an available feature. For this reason, I'll close and request it to be reopened.


After this has been done, I will review the pull request, and if anything needs to be 
changed or added, I will request it. Otherwise, I will accept the pull request and mark
the issue as completed and label it as `merged to devel`


[1] http://stackoverflow.com/questions/1783405/checkout-remote-git-branch
