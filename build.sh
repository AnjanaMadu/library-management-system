# Set the output directory for the binaries
OUTPUT_DIR="./output"

# Create the output directory if it doesn't exist
if [ ! -d $OUTPUT_DIR ]; then
  mkdir $OUTPUT_DIR
fi

# Set the name of the artifact
ARTIFACT_NAME="library-management-system"

# Set the list of target operating systems and architectures
TARGETS="windows/amd64 windows/386 linux/amd64 linux/386 linux/arm64 linux/arm darwin/amd64"

# Loop through the target operating systems and architectures
for TARGET in $TARGETS; do
  # Extract the operating system and architecture from the target
  OS=${TARGET%/*}
  ARCH=${TARGET#*/}

  # Set the output filename for the binary
  BIN_NAME=$ARTIFACT_NAME
  if [ $OS = "windows" ]; then
    BIN_NAME=$BIN_NAME.exe
  fi

  # Build the binary
  echo "Building $OS/$ARCH"
  GOOS=$OS GOARCH=$ARCH go build -o $BIN_NAME

  # Zip the binary and the public folder into a file named "lms-$OS-$ARCH.zip"
  7z a -tzip $OUTPUT_DIR/lms-$OS-$ARCH.zip $BIN_NAME public database.accdb

  # Delete the binary
  rm $BIN_NAME
done