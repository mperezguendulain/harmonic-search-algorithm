x1 = -10:0.5:10
y1 = x1
[x,y] = meshgrid(x1, y1); 

z = ((x.^2 - y.^2).*(sin(x+y)))+x+y
surfc(x,y,z) 