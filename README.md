# UE4GitDiff
Use Unreal Engine as a git difftool for comparing binary assets.

## Installation
Download the [latest release](https://github.com/mircearoata/UE4GitDiff/releases/latest) then run `UEGitDiff install` to register it as a difftool in the global .gitconfig (requires `git` to be under `PATH`).

## Usage
To run an asset diff, run `git difftool -t ue4 [rest of the arguments you would usually pass to git diff]`.
