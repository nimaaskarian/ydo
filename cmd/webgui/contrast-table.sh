IFS=$'\n' names=( $(grep -o 'color-.*:' "$(dirname "$0")/tailwind.css" | sed 's/color-//' | tr -d ':') )
IFS=$'\n' colors=( $(grep -o '#......' "$(dirname "$0")/tailwind.css") )

printf ","
for name in "${names[@]}"; do
  printf '"%s",' "$name"
done
printf "\n"

for i in "${!colors[@]}"; do
  printf '\"%s\",' "${names[i]}"
  for j in "${!colors[@]}"; do
    # github.com/nimaaskarian/contrast
    printf "%s, " "$(contrast "${colors[i]}" "${colors[j]}" -w)"
  done
  printf "\n"
done
