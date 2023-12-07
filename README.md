# Genesis

Evolution Simulator written in Go where creatures are governed by "concurrent" neural networks. 

This project is inspired by davidrandallmiller (https://github.com/davidrmiller/biosim4) written for practice and learning but there are some key differences the main being how creatures brains are actually simulated. In this simulator, a creature's brain output is calculated concurrently. This means that each input neuron is activated as a separate routine and a DFS search is conducted where each neighbor also starts another routine. This, in my opinion, better simulates an actual brain where different neurons fire at different times and can also be reactivated. This creates some very interesting brain structures and survival mechanisms. 

I also want to take this simulator several steps further by simulating climate to try to see migration patterns for example. I also want to add some "god" features so one can experiment with the organisms seeing how well they adapt to different scenarios. With the current simulator, I believe the organisms are very good at adapting leading to explosions in population early on. My goal now is to make it much harder to survive which is not a problem I foresaw having ;)
