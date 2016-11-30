args <- commandArgs(trailingOnly = TRUE)

csvContents <- read.csv(file=args[1], sep=",")
csvAsMatrx <- t(data.matrix(csvContents, rownames.force = NA))
print(csvAsMatrx)


graphName <- strsplit(args[1],"\\.")[[1]][1]
#fileName <- strcat(graphName, '.png')
png(file = paste(graphName, ".png", sep=""))

#barplot(csvAsMatrx)
barplot(csvAsMatrx, xlab="operation Sizes", ylab="time (ns)", main=graphName, names.arg=c("1 KB", "10 KB", "100 KB"))