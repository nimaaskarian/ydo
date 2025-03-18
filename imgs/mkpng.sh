for size in 48x48 32x32 16x16 192x192 180x180; do
  magick imgs/icon.png -resize $size cmd/webgui/static/imgs/icon-$size.png
done
