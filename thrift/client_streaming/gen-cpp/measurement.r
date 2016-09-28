# Read vals
data <- read.table("http://cs.wisc.edu/~riccardo/measurement.csv", header=T, sep=",") 



# Define colors to be used for cars, trucks, suvs
plot_colors <- c("blue", "red", "chartreuse4", "orange", "purple")

# Start PNG device driver to save output to figure.png
png(filename="/home/r/riccardo/figure.png", height=400, width=400, 
    bg="white")

# Graph autos using y axis that ranges from 0 to max_y.
# Turn off axes and annotations (axis labels) so we can 
# specify them ourself
bw = (data$TOTAL_BYTES * 1e-6 * 8)/(data$TIME_NS * 1e-9)

max_y = max(bw)

plot(x=data$TOTAL_BYTES * 0.001, y=bw, type="o", col=plot_colors[1], 
     ylim=c(0,max_y), ann=FALSE, xlab="RTT(ns)")


# Make x axis using N=2,4,8,16 labels
#axis(1, at=1:4, lab=c(2,4, 8, 16))

# Make y axis with horizontal labels that display ticks at 
# every 4 marks. 4*0:max_y is equivalent to c(0,4,8,12).
#axis(2, las=2, at=c(0, 2, 4, 6, 8, 10, 12))

# Create box around plot
box()


# Create a title with a red, bold/italic font
title(main="Thrift - Client Streaming", col.main="red", font.main=4)

# Label the x and y axes with dark green text
title(xlab= "KB Sent", col.lab=rgb(0,0.5,0))
title(ylab= "Bandwidth (Mbps)", col.lab=rgb(0,0.5,0))

# Create a legend at (1, max_y) that is slightly smaller 
# (cex) and uses the same line colors and points used by 
# the actual plots
#legend(1, max_y, c("Static (1)", "Static (2)", "Static (4)", "Static (8)", "Static (16)"), cex=0.8, col=plot_colors, 
#       pch=21:23, lty=1:3);

# Turn off device driver (to flush output to png)
dev.off()