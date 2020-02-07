 Release Process

The ITRIX-Edge Project is released on an as-needed basis. The process is as follows:

1. An issue is proposing a new release with a changelog since the last release
2. At least one of the [approvers](OWNERS_ALIASES) must approve this release
3. An approver creates [new release in GitHub](https://github.com/OP-Team/ITRIX-Edge) using a version and tag name like `vX.Y.Z` and attaching the release notes
4. An approver creates a release branch in the form `release-X.Y`
5. The `ITRIX-Edge_VERSION` variable is updated in `.gitlab-ci.yml`
6. The release issue is closed
7. The topic of the #ITRIX-Edge channel is updated with `vX.Y.Z is released! | ...`

## Major/minor releases, merge freezes and milestones

* ITRIX-Edge maintains one branch for major releases (vX.Y). Minor releases are available only as tags.

* Security patches and bugs might be backported.

* Fixes for major releases (vX.x.0) and minor releases (vX.Y.x) are delivered
  via maintenance releases (vX.Y.Z) and assigned to the corresponding open
  milestone (vX.Y). That milestone remains open for the major/minor releases
  support lifetime, which ends once the milestone closed. Then only a next major
  or minor release can be done.

* ITRIX-Edge major and minor releases are bound to the given ``ITRIX-Edge_version`` major/minor
  version numbers and other components' arbitrary versions, like etcd or network plugins.
  Older or newer versions are not supported and not tested for the given release.

* There is no unstable releases and no APIs, thus ITRIX-Edge doesn't follow
  [semver](http://semver.org/). Every version describes only a stable release.
  Breaking changes, if any introduced by changed defaults or non-contrib ansible roles'
  playbooks, shall be described in the release notes. Other breaking changes, if any in
  the contributed addons or bound versions of ITRIX-Edge and other components, are
  considered out of ITRIX-Edge scope and are up to the components' teams to deal with and
  document.

