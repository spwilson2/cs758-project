import csv

import numpy
import matplotlib
# Don't use X to display.
matplotlib.use('Agg')
import matplotlib.pyplot as pyplot

from constants import *
COLORS = [name for name, hex in matplotlib.colors.cnames.iteritems()]

def save_csv(results, file_):
    file_ = open(file_, 'w')
    w = csv.DictWriter(file_, fieldnames=results[0].keys())
    w.writeheader()
    w.writerows(results)
    file_.close()


def autolabel(rects, ax):
    # attach some text labels
    for rect in rects:
        height = rect.get_height()
        ax.text(rect.get_x() + rect.get_width()/2., 1.05*height,
                '%d' % int(height),
                ha='center', va='bottom')


def bar(results,
        file_,
        ylab=None,
        xlab=None,
        title=None):

    split_results = {}

    # Split results by type of op

    for result in results:
        op, length = result[Go.OP_KEY], result[Go.LENGTH_KEY]
        if op not in split_results:
            split_results[op] = []
        split_results[op].append(int(length))

    fig, ax = pyplot.subplots()

    width = 0.35 # width of the bars

    bars = []
    ops = []
    #x_loc_start = range(len(split_results.keys()))
    x_loc_start = [0]

    for op_num, (op, results) in enumerate(split_results.items()):

        # the x locations for op types
        x_loc = [val + width*op_num for val in x_loc_start]
        meanLength = numpy.mean(results)
        lengthStd  = numpy.std(results)
        bars.append(ax.bar(x_loc, meanLength, width, color=COLORS[op_num], yerr=lengthStd))
        ops.append(op)


    if ylab is not None:
        ax.set_ylabel(ylab)
    if title is not None:
        ax.set_title(title)
    if xlab is not None:
        if type(xlab) is not str:
            ax.set_xticks(x_loc_start)
        ax.set_xticklabels(xlab)

    for bar in bars:
        autolabel(bar, ax)

    ax.legend(bars, ops)

    pyplot.savefig(file_)
