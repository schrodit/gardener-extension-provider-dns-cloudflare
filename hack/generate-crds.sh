#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if ! command -v controller-gen >/dev/null 2>&1; then
  >&2 echo "controller-gen not available"
  exit 1
fi

output_dir="$(pwd)"
output_dir_temp="$(mktemp -d)"
file_name_prefix=""
crd_options=""
add_deletion_protection_label=false
add_keep_object_annotation=false
args=()

cleanup() {
  rm -rf "${output_dir_temp}"
}
trap cleanup EXIT

get_group_package() {
  case "$1" in
    "extensions.gardener.cloud")
      echo "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
      ;;
    "resources.gardener.cloud")
      echo "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
      ;;
    *)
      >&2 echo "unknown group $1"
      return 1
      ;;
  esac
}

add_metadata_mutations() {
  local file="$1"

  if "${add_deletion_protection_label}"; then
    if ! grep -q "clusters.extensions.gardener.cloud" "$file"; then
      perl -0pi -e 's/metadata:\n/metadata:\n  labels:\n    gardener.cloud\/deletion-protected: "true"\n/' "$file"
    fi
  fi

  if "${add_keep_object_annotation}"; then
    perl -0pi -e 's/annotations:\n/annotations:\n    resources.gardener.cloud\/keep-object: "true"\n/' "$file"
  fi
}

generate_group() {
  local group="$1"
  local package package_path sanitized_group_name pattern relevant_file file_name delete_file

  package="$(get_group_package "$group")"
  package_path="$(go list -mod=mod -f '{{ .Dir }}' "$package")"
  sanitized_group_name="${group%%_*}"

  controller-gen "crd${crd_options}" "paths=${package_path}" "output:crd:dir=${output_dir_temp}"

  relevant_file=""
  for crd in "${output_dir_temp}/${sanitized_group_name}"_*.yaml; do
    [ -e "$crd" ] || continue
    local crd_out="${output_dir}/${file_name_prefix}$(basename "$crd")"
    mv "$crd" "$crd_out"
    add_metadata_mutations "$crd_out"
    relevant_file="${relevant_file} $(basename "$crd_out")"
  done

  pattern=".*${group}_.*\\.yaml"
  for file in "${output_dir}"/*.yaml; do
    [ -e "$file" ] || continue
    file_name="$(basename "$file")"
    delete_file=true

    case " ${relevant_file} " in
      *" ${file_name} "*) delete_file=false ;;
    esac

    if [[ ! "$file_name" =~ $pattern ]]; then
      delete_file=false
    fi

    if "${delete_file}"; then
      rm "$file"
    fi
  done
}

while test $# -gt 0; do
  case "$1" in
    -p)
      file_name_prefix="$2"
      shift 2
      ;;
    -l)
      add_deletion_protection_label=true
      shift
      ;;
    -k)
      add_keep_object_annotation=true
      shift
      ;;
    --allow-dangerous-types)
      crd_options=":allowDangerousTypes=true"
      shift
      ;;
    *)
      args+=("$1")
      shift
      ;;
  esac
done

if [ "${#args[@]}" -eq 0 ]; then
  args=("extensions.gardener.cloud" "resources.gardener.cloud")
fi

for group in "${args[@]}"; do
  generate_group "$group"
done
