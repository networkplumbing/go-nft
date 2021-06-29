# Creating a New Release

## Changelog

- [ ] Update the CHANGELOG file in the project root directory.
  Use `git log  --oneline v0.1.0..HEAD` to get the changes since the last tag.
  Pick the most important changes, especially user facing ones.

  The following format should be used when adding an entry:

```
## [X.Y.Z] - YYYY-MM-DD
### Breaking Changes
 - ...

### New Features
 - ...

### Bug Fixes
 - ...
```

## Tagging

- [ ] Tag new release in git.
```bash
# Make sure your local git repo is sync with upstream.
# The whole version string should log like `v0.1.0`.
# For the commit message use the following format: `go-nft 0.1.0 release`.
git tag --sign v<version>
git push upstream --tags
```

- [ ] In case there is a need to remove a tag:
```bash
# Remove local tag
git tag -d <tag_name>

# Remove upstream tag
git push --delete upstream <tag_name>
```

## GitHub Release

- [ ] Visit [github draft release page][1].
- [ ] Make sure you are in `Release` tab.
- [ ] Choose the git tag just pushed.
- [ ] Set title with the following format: `Version 0.1.0 release`.
- [ ] The content should be copied from the `CHANGELOG` file.
- [ ] Click `Save draft` and seek for review.
- [ ] Click `Publish release` once approved.
