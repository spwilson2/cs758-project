import csv

import numpy
import matplotlib
# Don't use X to display.
matplotlib.use('Agg')
import matplotlib.pyplot as pyplot

from constants import *
COLORS = [name for name, hex in matplotlib.colors.cnames.items()]

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



def execution_bar(results, file_, ylab=None, title=None):

    '''Create a bar graph plotting all times on the same axis.'''
    split_results = {}

    # Split results by type of op
    for result in results:
        key = ""
        if result[Go.OP_KEY] == "Read":
            key = "rsize"
        elif result[Go.OP_KEY] == "Write":
            key = "wsize"
        else: 
            continue
        op_size, length = result[key], result[Go.LENGTH_KEY]
        if op_size not in split_results:
            split_results[op_size] = {}
        if result[Go.IO_TYPE_KEY] not in split_results[op_size]:
            split_results[op_size][result[Go.IO_TYPE_KEY]] = []
        split_results[op_size][result[Go.IO_TYPE_KEY]].append(int(length))

    processed_means = {}
    processed_std = {}
    for op_size in split_results.keys():
        for io_type in split_results[op_size]:
            if io_type not in processed_means:
                processed_means[io_type] = []
                processed_std[io_type] = []
            meanLength = numpy.mean(split_results[op_size][io_type])
            processed_means[io_type].append(meanLength)

            lengthStd = numpy.std(split_results[op_size][io_type])
            processed_std[io_type].append(lengthStd)

    N = len(split_results.keys())

    ind = numpy.arange(N)
    width = 0.25 # width of the bars

    fig, ax = pyplot.subplots()

    bars = []
    colors=['r','y','b','g']
    for op_num, io_type in enumerate(processed_means.keys()):
        bars.append(ax.bar((ind + (width * op_num)), processed_means[io_type], width, color=colors[op_num], yerr=processed_std[io_type]))

    if ylab is not None:
        ax.set_ylabel(ylab)
    if title is not None:
        ax.set_title(title)

    ax.set_xticks(ind + width)
    ax.set_xticklabels(split_results.keys())
    ax.set_ylim(bottom=0)

    ax.legend(bars, processed_means.keys(), loc=0)


    for bar in bars:
        autolabel(bar, ax)

    pyplot.savefig(file_)

def tracedata_bar(results, file_, ylab=None, title=None):

    '''Create a bar graph plotting all times on the same axis.'''
    split_results = {}

    # Split results by type of op
    for result in results:
        ind, length = result[sortParameter], result[Go.LENGTH_KEY]
        if ind not in split_results:
            split_results[ind] = []
        split_results[ind].append(int(length))

    fig, ax = pyplot.subplots()


    width = 0.35 # width of the bars

    bars = []
    ops = []
    datas = len(split_results[Go.READ_OP] if split_results[Go.READ_OP] else split_results[Go.WRITE_OP])
    x_loc_start = range(datas)
    for test_num, (ind, results) in enumerate(split_results.items()):
        if op != Go.READ_OP and op != Go.WRITE_OP:
            continue
        # the x locations for op types
        x_loc = [val + width*test_num for val in x_loc_start]
        meanLength = numpy.mean(results)
        lengthStd = numpy.std(results)
        bars.append(ax.bar(x_loc, meanLength, width, color=COLORS[test_num], yerr=lengthStd))
        ops.append(ind)

    pyplot.tick_params(axis='x', which='both', bottom='off', top='off', labelbottom='off')
    ax.set_ylim(bottom=0)
    if xlab is not None:
        ax.set_xlabel(xlab)
    if ylab is not None:
        ax.set_ylabel(ylab)
    if title is not None:
        ax.set_title(title)

    for bar in bars:
        autolabel(bar, ax)

    ax.legend(bars, ops, loc=4)

    pyplot.savefig(file_)