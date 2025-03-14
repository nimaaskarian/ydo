last_tag=$(git tag | tail -n 1)
git tag "$1" || last_tag=$(git tag | tail -n 2 | head -n 1)
make && \
gh release create "$1" ydo ydo.exe  --title "$1" --notes "**Full Changelog**: https://github.com/nimaaskarian/ydo/compare/$last_tag...$1" --repo nimaaskarian/ydo
