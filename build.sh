#!/bin/bash

# TODO implement port fallback logic 
set -eu
#set -x

trap 'code=$?; echo "[ERROR] (exit=$code) in ${FUNCNAME:-main} at line $LINENO" >&2; exit $code' ERR

# ===============================================================
#  colors
# ===============================================================

# Source - https://stackoverflow.com/a
# Posted by Shakiba Moshiri, modified by community. See post 'Timeline' for change history
# Retrieved 2025-11-07, License - CC BY-SA 4.0

# Reset
Coff='\033[0m'       # Text Reset

# Regular Colors
Black='\033[0;30m'        # Black
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        # White

# Bold
BBlack='\033[1;30m'       # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan
BWhite='\033[1;37m'       # White

# Underline
UBlack='\033[4;30m'       # Black
URed='\033[4;31m'         # Red
UGreen='\033[4;32m'       # Green
UYellow='\033[4;33m'      # Yellow
UBlue='\033[4;34m'        # Blue
UPurple='\033[4;35m'      # Purple
UCyan='\033[4;36m'        # Cyan
UWhite='\033[4;37m'       # White

# Background
On_Black='\033[40m'       # Black
On_Red='\033[41m'         # Red
On_Green='\033[42m'       # Green
On_Yellow='\033[43m'      # Yellow
On_Blue='\033[44m'        # Blue
On_Purple='\033[45m'      # Purple
On_Cyan='\033[46m'        # Cyan
On_White='\033[47m'       # White

# High Intensity
IBlack='\033[0;90m'       # Black
IRed='\033[0;91m'         # Red
IGreen='\033[0;92m'       # Green
IYellow='\033[0;93m'      # Yellow
IBlue='\033[0;94m'        # Blue
IPurple='\033[0;95m'      # Purple
ICyan='\033[0;96m'        # Cyan
IWhite='\033[0;97m'       # White

# Bold High Intensity
BIBlack='\033[1;90m'      # Black
BIRed='\033[1;91m'        # Red
BIGreen='\033[1;92m'      # Green
BIYellow='\033[1;93m'     # Yellow
BIBlue='\033[1;94m'       # Blue
BIPurple='\033[1;95m'     # Purple
BICyan='\033[1;96m'       # Cyan
BIWhite='\033[1;97m'      # White

# High Intensity backgrounds
On_IBlack='\033[0;100m'   # Black
On_IRed='\033[0;101m'     # Red
On_IGreen='\033[0;102m'   # Green
On_IYellow='\033[0;103m'  # Yellow
On_IBlue='\033[0;104m'    # Blue
On_IPurple='\033[0;105m'  # Purple
On_ICyan='\033[0;106m'    # Cyan
On_IWhite='\033[0;107m'   # White


# ===============================================================
#  Global variables
# ===============================================================
ROOT_DIR="sys"
BUILD_STATE_FILE=${BUILD_STATE_FILE:=".buildenv"}

FORMAT="${FORMAT:="0"}"
HELP=""
PLATFORM=""
CLEAN=""
BUILD=""
COVERAGE="${COVERAGE:="0"}"
BUILD_DIR=${BUILD_DIR:="bin"}
CC=""
CC_STD=${CC_STD:="-std=c11"}
UNPARSE=""
COMPLET_PATH_DIR=""
OBJ_DIR=""
EXE=""
LIST_LIBS="0"
LIB_OUTPUT_DIR=""
LIB_ORDER=""
EXECUTE="${EXECYTE:="0"}"


# ===============================================================
#  Persistent state management
# ===============================================================
load_env() {
	if [ ! -f "$BUILD_STATE_FILE" ]; then
		echo "[load_env] No previous build state found"
		return 0
	fi

	echo "[load_env] Loading build state from $BUILD_STATE_FILE"

	while IFS='=' read -r key value || [ -n "${key:-}" ]; do
		[ -z "${key:-}" ] && continue
		value="${value:-}"
		case "$key" in
			PLATFORM) [ -z "${PLATFORM:-}" ] && PLATFORM="$value" ;;
			BUILD)    [ -z "${BUILD:-}" ] && BUILD="$value" ;;
			CC)       [ -z "${CC:-}" ] && CC="$value" ;;
			CC_STD)   [ -z "${CC_STD:-}" ] && CC_STD="$value" ;;
			EXE)      [ -z "${EXE:-}" ] && EXE="$value" ;;
			CLEAN)    [ -z "${CLEAN:-}" ] && CLEAN="$value" ;;
			FORMAT)   [ -z "${FORMAT:-}" ] && FORMAT="$value" ;;
			LIB_ORDER) [ -z "${LIB_ORDER:-}" ] && LIB_ORDER="$value" ;;
			COVERAGE) [ -z "${COVERAGE:-}" ] && COVERAGE="$value" ;;
			EXECUTE) [ -z "${EXECUTE:-}" ] && EXECUTE="$value" ;;
		esac
	done < "$BUILD_STATE_FILE"
}

persist_var() {
	local key="$1"
	local value="$2"
	touch "$BUILD_STATE_FILE"
	if grep -q "^${key}=" "$BUILD_STATE_FILE" 2>/dev/null; then
		sed -i "s|^${key}=.*|${key}=${value}|" "$BUILD_STATE_FILE"
	else
		printf "%s=%s\n" "$key" "$value" >> "$BUILD_STATE_FILE"
	fi
}

newline() {
	# S'assurer que le fichier existe
	[ ! -e "$BUILD_STATE_FILE" ] && return 0

	# Vérifie si la dernière ligne est vide
	local lastline
	lastline=$(tail -n 1 "$BUILD_STATE_FILE" 2>/dev/null || echo "")

	# Si la dernière ligne n'est pas vide, ajoute un saut de ligne
	if [ -n "$lastline" ]; then
		printf "\n" >> "$BUILD_STATE_FILE"
	fi
}

save_env() {
	echo "[save_env] Writing build state to $BUILD_STATE_FILE"
	persist_var "PLATFORM" "$PLATFORM"
	persist_var "BUILD" "$BUILD"
	persist_var "COVERAGE" "$COVERAGE"
	persist_var "CLEAN" "$CLEAN"
	persist_var "CC" "$CC"
	persist_var "CC_STD" "$CC_STD"
	persist_var "EXE" "$EXE"
	persist_var "EXECUTE" "$EXECUTE"
	persist_var "FORMAT" "$FORMAT"
	persist_var "LIB_ORDER" "$LIB_ORDER"
	newline
}

reset() {
	if [ -z "${UNPARSE:-}" ]; then
		echo "  → skip reset"
		return 0
	fi
	echo "=== Resetting build state ==="
	rm -f "$BUILD_STATE_FILE" 2>/dev/null || true
}

# ===============================================================
#  Build helpers
# ===============================================================

detect_compiler() {
if [ -n "$CC" ]; then
	echo "[detect compiler] → skip"
	return 0
fi


if command -v clang >/dev/null; then
    CC=clang
elif command -v gcc >/dev/null; then
    CC=gcc
else
    echo "No compiler found."
    exit 1
fi

}

build_variables() {
	[ -z "$PLATFORM" ] && echo "[build_variables] no platform" && return 1
	[ -z "$BUILD" ] && echo "[build_variables] no build flag" && return 1
	[ -z "$EXE" ] && echo "[build_variables] no exe name" && return 1

	local path="${BUILD_DIR}/${BUILD}/${PLATFORM}"
	COMPLET_PATH_DIR="$path"
	OBJ_DIR="${path}/obj"
	LIB_OUTPUT_DIR="${COMPLET_PATH_DIR}/lib"
	EXE="${path}/${EXE}"
}

cflags() {
	# Flags communs (warnings stricts, sans optimisation initiale)
	local common_flags="-Wall -Wextra -Wpedantic \
		-Wshadow -Wcast-align -Wunused -Wold-style-definition \
		-Wmissing-prototypes -Wno-unused-parameter -Werror \
		-Wstrict-prototypes -Wpointer-arith -Wwrite-strings \
		-Wconversion -Wformat=2 -Wformat-security \
		-Wunreachable-code -Wundef -Wbad-function-cast \
		-Wdouble-promotion -Wmissing-include-dirs \
		-Winit-self -Wmissing-noreturn -fno-common \
		-fstack-protector-strong"

	local flags="$common_flags"

	if [ "$BUILD" = "debug" ]; then
		# Flags debug de base
		flags="$flags -g3 -O0 -DDEBUG -DDEBUG_MEMORY=1 -ftrapv"

		# Sanitizer (commun à Clang et GCC)
		local sanitizer_flags="-fsanitize=address,undefined -fno-omit-frame-pointer"
		flags="$flags $sanitizer_flags"

		# Coverage si activé
		if [ "$COVERAGE" = "1" ]; then
			local coverage_flags=""
			# Pour Clang (géré en priorité comme demandé)
			if [ "$CC" = "clang" ]; then
				coverage_flags="-fprofile-instr-generate -fcoverage-mapping"
				# Pour GCC (explicitement géré, comme demandé)
			elif [ "$CC" = "gcc" ]; then
				coverage_flags="-fprofile-arcs -ftest-coverage"
			else
				echo "Compilateur non supporté : $CC" >&2
				return 1
			fi
			flags="$flags $coverage_flags"
		fi
	elif [ "$BUILD" = "release" ]; then
		# Flags release
		flags="$flags -O2 -DNDEBUG -DDEBUG_MEMORY=0 -fomit-frame-pointer -march=native -D_FORTIFY_SOURCE=2"
	else
		echo "Valeur BUILD invalide : doit être 'debug' ou 'release'" >&2
		return 1
	fi

	echo "$flags"
}

include_flags() {
	local flags=""
	[ -d "$ROOT_DIR/include" ] && flags+=" -I$ROOT_DIR/include"
	[ -n "$PLATFORM" ] && [ -d "$PLATFORM" ] && flags+=" -I$PLATFORM"
	[ -d "$ROOT_DIR/include/sys" ] && flags+=" -I$ROOT_DIR/include/sys"
	echo "$flags"
}

link_flags() {
	echo "-lSDL2 -lm"
}

# ===============================================================
#  Modular build pipeline
# ===============================================================
discover_libs() {
	LIB_DIRS=()
	local base="${ROOT_DIR}/src"
	while IFS= read -r -d '' dir; do
		LIB_DIRS+=("$dir")
	done < <(find "$base" -maxdepth 1 -type d -name "lib*" -print0)
}

sources_lib() {
	local dir="$1"
	local list=()
	while IFS= read -r -d '' file; do
		list+=("$file")
	done < <(find "$dir" -type f -name "*.c" -print0)
	echo "${list[@]}"
}

compile_lib_sources() {
	local libname="$1"; shift
	local srcs=("$@")
	local outdir="$OBJ_DIR/$libname"
	mkdir -p "$outdir"

	local objs=()
	for src in "${srcs[@]}"; do
		local obj="$outdir/$(basename "$src" .c).o"
		# message utilisateur → stderr, jamais capturé
		echo "  CC $src" >&2
		"$CC" "$CC_STD" $(cflags) $(include_flags) -c "$src" -o "$obj"
		objs+=("$obj")
	done

	# Retourne seulement la liste des objets sur stdout
	printf "%s\n" "${objs[@]}"
}

build_lib() {
	local libname="$1"; shift
	local objs=("$@")
	local outlib="${LIB_OUTPUT_DIR}/${libname}.a"
	mkdir -p "${COMPLET_PATH_DIR}/lib"
	[ ${#objs[@]} -eq 0 ] && echo -e " $BRed  → no objects for $libname $Coff" && return 0
		ar rcs "$outlib" "${objs[@]}"
		echo -e "  $BGreen ✓  $Coff$outlib built"
	}

Build_all_libs() {
	discover_libs
	if [ ${#LIB_DIRS[@]} -eq 0 ]; then
		echo -e " $On_Red [build_all_libs] No libraries found under ${ROOT_DIR}/src $Coff"
		return 0
	fi
	for dir in "${LIB_DIRS[@]}"; do
		local name
		name=$(basename "$dir")
		echo -e  "$BPurple === Building library: $name === $Coff"
		local srcs=($(sources_lib "$dir"))
		local objs=($(compile_lib_sources "$name" "${srcs[@]}"))
		build_lib "$name" "${objs[@]}"
	done
}

list_libs() {
	if [ -z "${LIST_LIBS:-}" ] || [ "$LIST_LIBS" -eq 0 ]; then
		return 0
	fi

	discover_libs
	echo -e "📚  $Cyan Libraries detected in ${ROOT_DIR}/src: $Coff"
	if [ ${#LIB_DIRS[@]} -eq 0 ]; then
		echo "  → none found"
	else
		for dir in "${LIB_DIRS[@]}"; do
			echo "  - $(basename "$dir")"
		done
	fi
}

# ===============================================================
#  Platform backend
# ===============================================================
build_platform_backend() {
	[ -z "${PLATFORM:-}" ] && return 0
	[ ! -d "$PLATFORM" ] && return 0

	mkdir -p "$OBJ_DIR/$PLATFORM"
	echo -e "$BPurple === Building platform backend: $PLATFORM === $Coff"
	for src in "$PLATFORM"/*.c; do
		obj="$OBJ_DIR/$PLATFORM/$(basename "$src" .c).o"
		echo "  CC $src"
		"$CC" "$CC_STD" $(cflags) $(include_flags) -c "$src" -o "$obj"
	done
}


# ===============================================================
#  Discover compiled static archives (.a) and generate linker args
# ===============================================================
discover_archives() {
	local dir="${LIB_OUTPUT_DIR:-lib}"
	local libs="-L$dir "
	local found_any=0

	# Si LIB_ORDER est défini → suivre cet ordre
	if [ -n "$LIB_ORDER" ]; then
		for name in $LIB_ORDER; do
			local path="$dir/lib${name}.a"
			if [ -f "$path" ]; then
				libs="$libs -l${name} "
				found_any=1
			else
				echo -e "$BYellow [warn] expected lib${name}.a not found in $dir $Coff" >&2
			fi
		done
	else
		# sinon, fallback sur le scan automatique
		for file in "$dir"/lib*.a; do
			if [ -f "$file" ]; then
				local name=$(basename "$file")
				name=${name#lib}
				name=${name%.a}
				libs="$libs -l$name "
				found_any=1
			fi
		done
	fi

	if [ $found_any -eq 0 ]; then
		echo -e "$BYellow [warn] No static libraries (.a) found in '$dir' $Coff" >&2
	fi

	echo "$libs"
}


# ===============================================================
#  Linking phase
# ===============================================================
link_build() {
	local linker; linker="$(link_flags)"
	mkdir -p "$(dirname "$EXE")"
	echo -e "$BPurple === Linking final executable === $Coff"

	local lib_args
	lib_args=$(discover_archives)
	echo -e "${BGreen} [link] Using libraries:${Coff} $lib_args"

	local cmdsrc="$ROOT_DIR/src/cmd/engine/engine.c" # TODO variable
	local platform_objs=()
	[ -d "$OBJ_DIR/$PLATFORM" ] && platform_objs=("$OBJ_DIR/$PLATFORM"/*.o)

	"$CC" "$CC_STD" $(cflags) $(include_flags) \
		"$cmdsrc" "${platform_objs[@]}" \
		$lib_args  \
		$linker -o "$EXE"

	echo -e "{$BGreen}✓  Build OK $Coff → $Yellow  $EXE $Coff"
}

# ===============================================================
#  Cleaning & formatting
# ===============================================================
clean() {
	if [ -z "${CLEAN:-}" ] || [ "$CLEAN" -eq 0 ]; then
		echo "[clean] → skip"
		return 0
	fi
	echo "[clean] removing $BUILD_DIR and lib/"
	rm -rf "$BUILD_DIR" lib
}

format() {
	if [ -z "${FORMAT:-}" ] || [ "$FORMAT" -eq 0 ]; then
		echo "[format] → skip"
		return 0
	fi
	echo "[format] Running clang-format pass..."
	mapfile -t files < <(find . -type f \( -name '*.c' -o -name '*.h' \))
	[ ${#files[@]} -eq 0 ] && echo "  → no source files found" && return 0
		clang-format -i "${files[@]}"
		echo -e "  $BGreen ✓  $Coff clang -format applied"
	}

# ===============================================================
#  Coverage / Analyse
# ===============================================================

analyse_verify() {
    # Vérifie la présence de llvm-profdata et llvm-cov
    if ! command -v llvm-profdata >/dev/null 2>&1 || ! command -v llvm-cov >/dev/null 2>&1; then
        echo "0"
    else
        echo "1"
    fi
}

analyse_clean() {
    echo -e "${BPurple}🧹 Cleaning previous coverage data...${Coff}"

    find . \( \
        -name '*.gcda' -o \
        -name '*.gcno' -o \
        -name 'coverage.profraw' -o \
        -name 'coverage.profdata' \
    \) -type f -print -delete 2>/dev/null || true

    [ -d coverage_report ] && rm -rf coverage_report

	echo -e "${BPurple}🧹 Cleaning previous coverage data...${Coff}"
}

analyze_coverage_gcc() {
    echo -e "${BPurple}📊 Running GCC coverage analysis...${Coff}"
    gcov "$EXE"
}

analyze_coverage_llvm() {
    echo -e "${BPurple}🔨 Building with coverage instrumentation...${Coff}"

    local report_dir="${COMPLET_PATH_DIR}/coverage"
    mkdir -p "$report_dir"

    echo -e "${BPurple}🚀 Running program to generate coverage data...${Coff}"
    LLVM_PROFILE_FILE="${report_dir}/coverage.profraw" "./$EXE"

    echo -e "${BPurple}📊 Processing coverage data...${Coff}"
    llvm-profdata merge -sparse "${report_dir}/coverage.profraw" -o "${report_dir}/coverage.profdata"

    echo -e "${BPurple}📋 Generating coverage report (console)...${Coff}"
    llvm-cov show "./$EXE" \
        -instr-profile="${report_dir}/coverage.profdata" \
        --show-line-counts-or-regions \
        --show-expansions \
        "$ROOT_DIR/src" \
        > "${report_dir}/coverage_report.txt"

    echo -e "${BPurple}📝 Generating coverage summary...${Coff}"
    llvm-cov report "./$EXE" \
        -instr-profile="${report_dir}/coverage.profdata" \
        "$ROOT_DIR/src" \
        > "${report_dir}/coverage_summary.txt"

    echo -e "${BGreen}✓ Coverage reports generated in: ${report_dir}${Coff}"
    echo -e "  ├─ coverage.profraw"
    echo -e "  ├─ coverage.profdata"
    echo -e "  ├─ coverage_report.txt"
    echo -e "  └─ coverage_summary.txt"
}

analyze_coverage_html() {
    analyze_coverage_llvm  # génère déjà le .profdata et les rapports texte

    local report_dir="${COMPLET_PATH_DIR}/coverage"
    echo -e "${BPurple}🌐 Generating HTML coverage report...${Coff}"

    llvm-cov show "./$EXE" \
        -instr-profile="${report_dir}/coverage.profdata" \
        -format=html \
        -output-dir="${report_dir}/html" \
        -show-line-counts-or-regions \
        -show-expansions \
        -show-instantiations \
        "$ROOT_DIR/src"

    echo -e "${BGreen}✨ HTML report → ${report_dir}/html/index.html${Coff}"

    if command -v xdg-open >/dev/null 2>&1; then
        xdg-open "${report_dir}/html/index.html"
    elif command -v open >/dev/null 2>&1; then
        open "${report_dir}/html/index.html"
    else
        echo "🔔 Please open ${report_dir}/html/index.html manually"
    fi
}

analyse() {
    if [ "$BUILD" = "release" ] || [ "$COVERAGE" = "0" ] || [ "$(analyse_verify)" = "0" ]; then
        return 0
    fi

    echo -e "${BPurple}=== Running coverage analysis ===${Coff}"
    analyse_clean
	
    case "$CC" in
        "clang")
			echo -e "${BGreen} detect clang ${Coff}"
            analyze_coverage_html
            ;;
        "gcc")
			echo -e "${BGreen} detect gcc ${Coff}"
            analyze_coverage_gcc
            ;;
        *)
            echo -e "${BYellow}⚠️ Coverage analysis not supported for compiler: $CC${Coff}"
            ;;
    esac
}

execute(){
	if [ "$EXECUTE" == "0" ] || [ ! -e "$EXE" ]; then
		return 0
	fi

	$EXE
}


# ===============================================================
#  Help
# ===============================================================
helper() {
	if [ -z "$HELP" ]; then return 0; fi

	case "$HELP" in
		all|"")
			echo
			echo -e "$BPurple Usage: $Coff $Yellow ./build.sh [options] $Coff"
			echo
			echo -e "$BPurple Options: $Coff"
			printf " $Yellow  %-18s $Coff %s\n" "--platform <name>" "Target backend (e.g. sdl, x11, win32)"
			printf " $Yellow  %-18s $Coff %s\n" "--cc <compiler>" "Choose compiler (clang, gcc)"
			printf " $Yellow  %-18s $Coff %s\n" "--std <version>" "Set C standard (default: c11)"
			printf " $Yellow  %-18s $Coff %s\n" "--build <mode>" "Build type: release or debug"
			printf " $Yellow  %-18s $Coff %s\n" "--output, -o <name>" "Executable name"
			printf " $Yellow  %-18s $Coff %s\n" "--clean, -c [0|1]" "Clean build directories"
			printf " $Yellow  %-18s $Coff %s\n" "--format, -f [0|1]" "Apply clang-format"
			printf " $Yellow  %-18s $Coff %s\n" "--list-libs" "List detected libraries (no build)"
			printf " $Yellow  %-18s $Coff %s\n" "--reset, -u" "Reset build environment"
			printf " $Yellow  %-18s $Coff %s\n" "--help" "Display this help" 
			printf " $Yellow  %-18s $Coff %s\n" "--lib-order" "set link library order"
			printf " $Yellow  %-18s $Coff %s\n" "--coverage" "add coverage flags and automatically execute coverage tools"
			echo
			echo -e "$BPurple Build process overview: $Coff"
			echo -e "  1. Discover all libs under $Green sys/src/lib* $Coff"
			echo -e "  2. Compile each lib into  $Green lib/libXXX.a $Coff"
			echo -e "  3. Compile backend sources if $Yellow --platform $Coff is set"
			echo -e "  4. Link all $Green .a $Coff + backend +  $Green cmd/engine.c $Coff  into the final binary"
			echo
			echo -e "$BPurple Examples: $Coff"
			echo -e "  $Yellow ./build.sh --platform sdl --cc clang --build release --output engine $Coff"
			echo -e "  $Yellow ./build.sh --list-libs $Coff"
			echo
			;;
		*)
			echo -e "$Red Unknown help topic: '$HELP' $Coff"
			;;
	esac
	exit 0
}

# ===============================================================
#  Argument parsing
# ===============================================================
parse_clean() { case "$1" in ""|1|true|yes|y) echo "1";; *) echo "0";; esac; }
parse_coverage() { case "$1" in ""|1|true|yes|y) echo "1";; *) echo "0";; esac; }
parse_execute() { case "$1" in ""|1|true|yes|y) echo "1";; *) echo "0";; esac; }
parse_build() { case "$1" in debug|0|false|no|n) echo "debug";; *) echo "release";; esac; }
parse_format() { case "$1" in ""|1|true|yes|y) echo "1";; *) echo "0";; esac; }
parse_std() { [ -z "$1" ] && echo "$CC_STD" || echo "-std=$1"; }

parse_lib_order() {
	# On garde la chaîne telle quelle, mais on peut la normaliser plus tard si besoin
	echo "$1"
}

parse() {
	while [ $# -gt 0 ]; do
		case "$1" in
			platform|--platform) shift; PLATFORM="$1"; shift ;;
			cc|--cc) shift; CC="$1"; shift ;;
			output|--output|-o) shift; EXE="$1"; shift ;;
			clean|--clean|-c) shift; CLEAN=$(parse_clean "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			build|--build|-b) shift; BUILD=$(parse_build "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			format|--format|-f) shift; FORMAT=$(parse_format "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			reset|--reset|-u) shift; UNPARSE="1" ;;
			std|--std) shift; CC_STD=$(parse_std "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			help|--help) shift; HELP="${1:-all}"; [ $# -gt 0 ] && shift || true ;;
			list-libs|--list-libs) LIST_LIBS="1"; shift ;;
			lib-order|--lib-order) shift; LIB_ORDER=$(parse_lib_order "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			coverage|--coverage) shift; COVERAGE=$(parse_coverage "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			execute |--execute|e|-e|--exec) shift; EXECUTE=$(parse_execute "${1:-}"); [ $# -gt 0 ] && shift || true ;;
			*) echo "Unknown parameter: $1"; exit 1 ;;
		esac
	done
}

# ===============================================================
#  Main
# ===============================================================
main() {
	load_env
	parse "$@"
	detect_compiler
	save_env
	helper
	reset
	format
	clean
	build_variables
	list_libs
	build_all_libs
	build_platform_backend
	link_build
	analyse
	execute
}

main "$@"
