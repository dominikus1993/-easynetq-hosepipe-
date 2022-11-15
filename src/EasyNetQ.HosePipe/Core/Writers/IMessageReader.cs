using System.Threading.Channels;

using EasyNetQ.HosePipe.Types;

namespace EasyNetQ.HosePipe.Core.Writers;

public readonly record struct Queue(string Name);

public interface IMessageReader
{
    IAsyncEnumerable<HosepipeMessage> ReadMessages(Queue queue, CancellationToken cancellationToken);
}