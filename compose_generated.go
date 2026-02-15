package optic

// Compose3 returns an [Optic] composed of the 3 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose3[I1 any,I2 any,I3 any,S any,T any,A any,B any,C any,D any,E any,F any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3])Optic[I3,S,T,E,F,CompositionTree[CompositionTree[RET1,RET2],RET3],CompositionTree[CompositionTree[RW1,RW2],RW3],CompositionTree[CompositionTree[DIR1,DIR2],DIR3],CompositionTree[CompositionTree[ERR1,ERR2],ERR3]]{
return Compose(Compose(o1,o2),o3)
}
// Compose4 returns an [Optic] composed of the 4 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose4[I1 any,I2 any,I3 any,I4 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4])Optic[I4,S,T,G,H,CompositionTree[CompositionTree[RET1,RET2],CompositionTree[RET3,RET4]],CompositionTree[CompositionTree[RW1,RW2],CompositionTree[RW3,RW4]],CompositionTree[CompositionTree[DIR1,DIR2],CompositionTree[DIR3,DIR4]],CompositionTree[CompositionTree[ERR1,ERR2],CompositionTree[ERR3,ERR4]]]{
return Compose(Compose(o1,o2),Compose(o3,o4))
}
// Compose5 returns an [Optic] composed of the 5 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose5[I1 any,I2 any,I3 any,I4 any,I5 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,J any,K any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any,RET5 any,RW5 any,DIR5 any,ERR5 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4],o5 Optic[I5,G,H,J,K,RET5,RW5,DIR5,ERR5])Optic[I5,S,T,J,K,CompositionTree[CompositionTree[CompositionTree[RET1,RET2],RET3],CompositionTree[RET4,RET5]],CompositionTree[CompositionTree[CompositionTree[RW1,RW2],RW3],CompositionTree[RW4,RW5]],CompositionTree[CompositionTree[CompositionTree[DIR1,DIR2],DIR3],CompositionTree[DIR4,DIR5]],CompositionTree[CompositionTree[CompositionTree[ERR1,ERR2],ERR3],CompositionTree[ERR4,ERR5]]]{
return Compose(Compose(Compose(o1,o2),o3),Compose(o4,o5))
}
// Compose6 returns an [Optic] composed of the 6 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose6[I1 any,I2 any,I3 any,I4 any,I5 any,I6 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,J any,K any,L any,M any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any,RET5 any,RW5 any,DIR5 any,ERR5 any,RET6 any,RW6 any,DIR6 any,ERR6 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4],o5 Optic[I5,G,H,J,K,RET5,RW5,DIR5,ERR5],o6 Optic[I6,J,K,L,M,RET6,RW6,DIR6,ERR6])Optic[I6,S,T,L,M,CompositionTree[CompositionTree[CompositionTree[RET1,RET2],RET3],CompositionTree[CompositionTree[RET4,RET5],RET6]],CompositionTree[CompositionTree[CompositionTree[RW1,RW2],RW3],CompositionTree[CompositionTree[RW4,RW5],RW6]],CompositionTree[CompositionTree[CompositionTree[DIR1,DIR2],DIR3],CompositionTree[CompositionTree[DIR4,DIR5],DIR6]],CompositionTree[CompositionTree[CompositionTree[ERR1,ERR2],ERR3],CompositionTree[CompositionTree[ERR4,ERR5],ERR6]]]{
return Compose(Compose(Compose(o1,o2),o3),Compose(Compose(o4,o5),o6))
}
// Compose7 returns an [Optic] composed of the 7 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose7[I1 any,I2 any,I3 any,I4 any,I5 any,I6 any,I7 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,J any,K any,L any,M any,N any,O any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any,RET5 any,RW5 any,DIR5 any,ERR5 any,RET6 any,RW6 any,DIR6 any,ERR6 any,RET7 any,RW7 any,DIR7 any,ERR7 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4],o5 Optic[I5,G,H,J,K,RET5,RW5,DIR5,ERR5],o6 Optic[I6,J,K,L,M,RET6,RW6,DIR6,ERR6],o7 Optic[I7,L,M,N,O,RET7,RW7,DIR7,ERR7])Optic[I7,S,T,N,O,CompositionTree[CompositionTree[CompositionTree[RET1,RET2],CompositionTree[RET3,RET4]],CompositionTree[CompositionTree[RET5,RET6],RET7]],CompositionTree[CompositionTree[CompositionTree[RW1,RW2],CompositionTree[RW3,RW4]],CompositionTree[CompositionTree[RW5,RW6],RW7]],CompositionTree[CompositionTree[CompositionTree[DIR1,DIR2],CompositionTree[DIR3,DIR4]],CompositionTree[CompositionTree[DIR5,DIR6],DIR7]],CompositionTree[CompositionTree[CompositionTree[ERR1,ERR2],CompositionTree[ERR3,ERR4]],CompositionTree[CompositionTree[ERR5,ERR6],ERR7]]]{
return Compose(Compose(Compose(o1,o2),Compose(o3,o4)),Compose(Compose(o5,o6),o7))
}
// Compose8 returns an [Optic] composed of the 8 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose8[I1 any,I2 any,I3 any,I4 any,I5 any,I6 any,I7 any,I8 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,J any,K any,L any,M any,N any,O any,P any,Q any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any,RET5 any,RW5 any,DIR5 any,ERR5 any,RET6 any,RW6 any,DIR6 any,ERR6 any,RET7 any,RW7 any,DIR7 any,ERR7 any,RET8 any,RW8 any,DIR8 any,ERR8 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4],o5 Optic[I5,G,H,J,K,RET5,RW5,DIR5,ERR5],o6 Optic[I6,J,K,L,M,RET6,RW6,DIR6,ERR6],o7 Optic[I7,L,M,N,O,RET7,RW7,DIR7,ERR7],o8 Optic[I8,N,O,P,Q,RET8,RW8,DIR8,ERR8])Optic[I8,S,T,P,Q,CompositionTree[CompositionTree[CompositionTree[RET1,RET2],CompositionTree[RET3,RET4]],CompositionTree[CompositionTree[RET5,RET6],CompositionTree[RET7,RET8]]],CompositionTree[CompositionTree[CompositionTree[RW1,RW2],CompositionTree[RW3,RW4]],CompositionTree[CompositionTree[RW5,RW6],CompositionTree[RW7,RW8]]],CompositionTree[CompositionTree[CompositionTree[DIR1,DIR2],CompositionTree[DIR3,DIR4]],CompositionTree[CompositionTree[DIR5,DIR6],CompositionTree[DIR7,DIR8]]],CompositionTree[CompositionTree[CompositionTree[ERR1,ERR2],CompositionTree[ERR3,ERR4]],CompositionTree[CompositionTree[ERR5,ERR6],CompositionTree[ERR7,ERR8]]]]{
return Compose(Compose(Compose(o1,o2),Compose(o3,o4)),Compose(Compose(o5,o6),Compose(o7,o8)))
}
// Compose9 returns an [Optic] composed of the 9 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose9[I1 any,I2 any,I3 any,I4 any,I5 any,I6 any,I7 any,I8 any,I9 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,J any,K any,L any,M any,N any,O any,P any,Q any,R any,U any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any,RET5 any,RW5 any,DIR5 any,ERR5 any,RET6 any,RW6 any,DIR6 any,ERR6 any,RET7 any,RW7 any,DIR7 any,ERR7 any,RET8 any,RW8 any,DIR8 any,ERR8 any,RET9 any,RW9 any,DIR9 any,ERR9 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4],o5 Optic[I5,G,H,J,K,RET5,RW5,DIR5,ERR5],o6 Optic[I6,J,K,L,M,RET6,RW6,DIR6,ERR6],o7 Optic[I7,L,M,N,O,RET7,RW7,DIR7,ERR7],o8 Optic[I8,N,O,P,Q,RET8,RW8,DIR8,ERR8],o9 Optic[I9,P,Q,R,U,RET9,RW9,DIR9,ERR9])Optic[I9,S,T,R,U,CompositionTree[CompositionTree[CompositionTree[CompositionTree[RET1,RET2],RET3],CompositionTree[RET4,RET5]],CompositionTree[CompositionTree[RET6,RET7],CompositionTree[RET8,RET9]]],CompositionTree[CompositionTree[CompositionTree[CompositionTree[RW1,RW2],RW3],CompositionTree[RW4,RW5]],CompositionTree[CompositionTree[RW6,RW7],CompositionTree[RW8,RW9]]],CompositionTree[CompositionTree[CompositionTree[CompositionTree[DIR1,DIR2],DIR3],CompositionTree[DIR4,DIR5]],CompositionTree[CompositionTree[DIR6,DIR7],CompositionTree[DIR8,DIR9]]],CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR1,ERR2],ERR3],CompositionTree[ERR4,ERR5]],CompositionTree[CompositionTree[ERR6,ERR7],CompositionTree[ERR8,ERR9]]]]{
return Compose(Compose(Compose(Compose(o1,o2),o3),Compose(o4,o5)),Compose(Compose(o6,o7),Compose(o8,o9)))
}
// Compose10 returns an [Optic] composed of the 10 input optics.
//
// Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.
//
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [Compose] for a version that takes 2 parameters.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
func Compose10[I1 any,I2 any,I3 any,I4 any,I5 any,I6 any,I7 any,I8 any,I9 any,I10 any,S any,T any,A any,B any,C any,D any,E any,F any,G any,H any,J any,K any,L any,M any,N any,O any,P any,Q any,R any,U any,V any,W any,RET1 any,RW1 any,DIR1 any,ERR1 any,RET2 any,RW2 any,DIR2 any,ERR2 any,RET3 any,RW3 any,DIR3 any,ERR3 any,RET4 any,RW4 any,DIR4 any,ERR4 any,RET5 any,RW5 any,DIR5 any,ERR5 any,RET6 any,RW6 any,DIR6 any,ERR6 any,RET7 any,RW7 any,DIR7 any,ERR7 any,RET8 any,RW8 any,DIR8 any,ERR8 any,RET9 any,RW9 any,DIR9 any,ERR9 any,RET10 any,RW10 any,DIR10 any,ERR10 any](o1 Optic[I1,S,T,A,B,RET1,RW1,DIR1,ERR1],o2 Optic[I2,A,B,C,D,RET2,RW2,DIR2,ERR2],o3 Optic[I3,C,D,E,F,RET3,RW3,DIR3,ERR3],o4 Optic[I4,E,F,G,H,RET4,RW4,DIR4,ERR4],o5 Optic[I5,G,H,J,K,RET5,RW5,DIR5,ERR5],o6 Optic[I6,J,K,L,M,RET6,RW6,DIR6,ERR6],o7 Optic[I7,L,M,N,O,RET7,RW7,DIR7,ERR7],o8 Optic[I8,N,O,P,Q,RET8,RW8,DIR8,ERR8],o9 Optic[I9,P,Q,R,U,RET9,RW9,DIR9,ERR9],o10 Optic[I10,R,U,V,W,RET10,RW10,DIR10,ERR10])Optic[I10,S,T,V,W,CompositionTree[CompositionTree[CompositionTree[CompositionTree[RET1,RET2],RET3],CompositionTree[RET4,RET5]],CompositionTree[CompositionTree[CompositionTree[RET6,RET7],RET8],CompositionTree[RET9,RET10]]],CompositionTree[CompositionTree[CompositionTree[CompositionTree[RW1,RW2],RW3],CompositionTree[RW4,RW5]],CompositionTree[CompositionTree[CompositionTree[RW6,RW7],RW8],CompositionTree[RW9,RW10]]],CompositionTree[CompositionTree[CompositionTree[CompositionTree[DIR1,DIR2],DIR3],CompositionTree[DIR4,DIR5]],CompositionTree[CompositionTree[CompositionTree[DIR6,DIR7],DIR8],CompositionTree[DIR9,DIR10]]],CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR1,ERR2],ERR3],CompositionTree[ERR4,ERR5]],CompositionTree[CompositionTree[CompositionTree[ERR6,ERR7],ERR8],CompositionTree[ERR9,ERR10]]]]{
return Compose(Compose(Compose(Compose(o1,o2),o3),Compose(o4,o5)),Compose(Compose(Compose(o6,o7),o8),Compose(o9,o10)))
}
