namespace EasyNetQ.HosePipe.Types;

public sealed record HosepipeMessage(ReadOnlyMemory<byte> Body, MessageProperties Properties, MessageReceivedInfo Info);