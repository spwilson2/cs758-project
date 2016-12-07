import csv

#import pdb
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

    if len(split_results.keys()) == 0:
        return

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

    N = len(processed_means[processed_means.keys()[0]])

    ind = numpy.arange(N)
    width = 0.35 # width of the bars

    fig, ax = pyplot.subplots()

    bars = []
    colors=['r','y','c','w','m']
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
        key = ""
        if "rsize" in result.keys():
            key = "rsize"
        elif "wsize" in result.keys():
            key = "wsize"

        # we only want trace data which only occurs on non blocking operations
        if result[Go.OP_KEY] == "Read" or result[Go.OP_KEY] == "Write" or result[Go.IO_TYPE_KEY] == "blocking":
            continue

        op_size, length = result[key], result[Go.LENGTH_KEY]
        if op_size not in split_results:
            split_results[op_size] = {}
        if result[Go.OP_KEY] not in split_results[op_size]:
            split_results[op_size][result[Go.OP_KEY]] = []
        split_results[op_size][result[Go.OP_KEY]].append(int(length))

    if len(split_results.keys()) == 0:
        return

    processed_means = {}
    processed_std = {}
    for op_size in split_results.keys():
        for sched_overhead in split_results[op_size]:
            if sched_overhead not in processed_means:
                processed_means[sched_overhead] = []
                processed_std[sched_overhead] = []
            meanLength = numpy.mean(split_results[op_size][sched_overhead])
            processed_means[sched_overhead].append(meanLength)

            lengthStd = numpy.std(split_results[op_size][sched_overhead])
            processed_std[sched_overhead].append(lengthStd)

    N = len(processed_means[processed_means.keys()[0]])

    ind = numpy.arange(N)
    width = 0.15 # width of the bars

    fig, ax = pyplot.subplots(figsize=(20,8))

    bars = []
    colors=['r','y','c','w','m']

    for op_num, op_type in enumerate(processed_means.keys()):
        bars.append(ax.bar((ind + (width * op_num)), processed_means[op_type], width, color=colors[op_num], yerr=processed_std[op_type]))

    if ylab is not None:
        ax.set_ylabel(ylab)
    if title is not None:
        ax.set_title(title)

    ax.set_xticks(ind + (width * (N / 2)))
    ax.set_xticklabels(split_results.keys())
    ax.set_ylim(bottom=0)

    ax.legend(bars, processed_means.keys(), loc=0)


    for bar in bars:
        autolabel(bar, ax)

    pyplot.savefig(file_)

def flat_bar(results,
        file_,
        ylab=None,
        xlab=None,
        title=None):

    '''Create a bar graph plotting all times on the same axis.'''
    split_results = {}

    # Split results by type of op

    #pdb.set_trace()
    for result in results:
        op, length = result[Go.OP_KEY], result[Go.LENGTH_KEY]
        if op not in split_results:
            split_results[op] = []
        split_results[op].append(int(length))

    fig, ax = pyplot.subplots()


    width = 0.35 # width of the bars
    #width = 1 # width of the bars

    bars = []
    ops = []
    datas = len(split_results[Go.READ_OP] if Go.READ_OP in split_results else split_results[Go.WRITE_OP])
    x_loc_start = range(datas)
    #width =1/datas # width of the bars

    for op_num, (op, results) in enumerate(split_results.items()):
        if op != Go.READ_OP and op != Go.WRITE_OP:
            continue

        # the x locations for op types
        x_loc = [val + width*op_num for val in x_loc_start]
        lengthStd  = numpy.std(results)
        bars.append(ax.bar(x_loc, results, width, color=COLORS[op_num]))
        ops.append(op)


    pyplot.tick_params(axis='x', which='both', bottom='off', top='off', labelbottom='off')
    ax.set_ylim(bottom=0)
    if ylab is not None:
        ax.set_ylabel(ylab)
    if title is not None:
        ax.set_title(title)
    if xlab is not None:
        if type(xlab) is not str:
            ax.set_xticks(x_loc_start)
        ax.set_xticklabels(xlab)

    #for bar in bars:
    #    autolabel(bar, ax)

    ax.legend(bars, ops)

    pyplot.savefig(file_)
