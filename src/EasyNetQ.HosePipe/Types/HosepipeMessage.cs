namespace EasyNetQ.HosePipe.Types;

public sealed record HosepipeMessage(string Body, MessageProperties Properties, MessageReceivedInfo Info);