bin := net_utils
src_dir := ./utils/

${bin}: ${src_dir}/*.go
	cd ${src_dir} && go build -o ../$@ .
