# CoFaaS - Automatic Transformation-based Consolidation of Serverless Functions

## Introduction
An attractive property of serverless architectures is that they enable applications to be composed by multiple individual functions that may be written in different languages. The main challenge of these architectures is that inter-function communication comes with significant overheads. To address this, CoFaaS proposes to automatically consolidate modular serverless applications onto a single runtime to achieve inter-function communication latencies approaching those of a function call within a monolithic application. To do this, CoFaaS relies on the well-defined external APIs often used by FaaS functions. These well-defined APIs are declared in dedicated Interface Description Languages (IDLs) such as Protobuf for gRPC. CoFaaS transforms the IDL descriptions to WebAssembly Interface Types (WIT), the IDL used by the WebAssemby component model. Thereby, FaaS functions can be compiled to a WebAssembly component that can be Co-located on a single WebAssembly runtime.

## Dependencies
To use this project, make sure the following binaries are available in `$PATH`



## Repository overview


## Project status
CoFaaS is an academic research proposal supported by a
proof-of-concept implementation. This work, the paper and source code,
are provided in the hope that the implemented techniques will be
useful to other similar projects. Following the proof-of-concept
nature of the work, CoFaaS is currently very limited and no further
development is planned.

## Referenceing CoFaaS
If you use or derive from CoFaaS in your academic research, please cite our workshop accepted to the 2nd SESAME workshop at EuroSys 2024:

```
@inproceedings{asheim24_cofaas,
  author =       {Asheim, Truls and Jahre, Magnus and Kumar, Rakesh},
  title =        {CoFaaS: Automatic Transformation-based Consolidation
                  of Serverless Functions},
  booktitle =    {Proceedings of the 2nd Workshop on SErverless
                  Systems, Applications and MEthodologies},
  year =         2024,
  doi =          {10.1145/3642977.3652093},
  publisher =    {Association for Computing Machinery},
}
```

## License
