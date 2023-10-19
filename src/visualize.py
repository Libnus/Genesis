#!/usr/bin/env python

import pandas as pd
import networkx as nx
import matplotlib.pyplot as plt
from pyvis.network import Network
import scipy

edge_list = pd.read_csv("neural_output.brain")

G = nx.from_pandas_edgelist(edge_list,
                            source="source",
                            target="target",
                            edge_attr = ["value"],
                            create_using=nx.DiGraph())

nodeColors = {}
for i, row in edge_list.iterrows():
    nodeColors[row['source']] = row['color1']
    nodeColors[row['target']] = row['color2']

# plt.figure(figsize=(10,10))
# pos = nx.spring_layout(G, seed=13648)
# nx.draw(G, with_labels=True, node_color="red", edge_cmap=plt.cm.Blues, pos=pos)
nx.set_node_attributes(G, 50, 'size')
nx.set_node_attributes(G, nodeColors, name="color")
# plt.show()

net = Network(directed=True, width="1400px", height="900px", bgcolor="#222222", font_color="white")
# net.repulsion()
net.show_buttons(filter_=["physics"])
net.barnes_hut()
net.from_nx(G)

net.show("neural.html", notebook=False)
