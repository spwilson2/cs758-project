args <- commandArgs(trailingOnly = TRUE)

csvContents <- read.csv(file=args[1], sep=",", header=FALSE)
csvAsMatrx <- data.matrix(csvContents, rownames.force = NA)
print(csvAsMatrx)


graphName <- strsplit(strsplit(args[1],"\\.")[[1]][2], "\\/")[[1]][3]
png(file = paste("graphs/", graphName, ".png", sep=""))

#barplot(csvAsMatrx)
barplot(csvAsMatrx, xlab="operation Sizes", ylab="time (ns)", main=graphName, names.arg=c("1 KB", "10 KB", "100 KB", "1 MB"))