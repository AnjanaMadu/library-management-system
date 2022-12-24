# Set the output directory for the binaries
OUTPUT_DIR="./build"

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
  GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_DIR/$BIN_NAME

  # Zip the binary and the public folder into a file named "lms-$OS-$ARCH.zip"
  zip $OUTPUT_DIR/lms-$OS-$ARCH.zip $OUTPUT_DIR/$BIN_NAME public database.accdb

  # Delete the binary
  rm $OUTPUT_DIR/$BIN_NAME
done