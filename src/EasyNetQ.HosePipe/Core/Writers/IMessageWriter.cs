using EasyNetQ.HosePipe.Types;

namespace EasyNetQ.HosePipe.Core.Writers;

public interface IMessageWriter
{
    Task Write(IEnumerable<HosepipeMessage> messages, QueueParameters parameters, CancellationToken cancellationToken);
}