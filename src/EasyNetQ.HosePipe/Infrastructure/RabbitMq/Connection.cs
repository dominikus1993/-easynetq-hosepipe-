using RabbitMQ.Client;

namespace EasyNetQ.HosePipe.Infrastructure.RabbitMq;

public static class Connection
{
    public static IConnection FromUrl(string? url)
    {
        ArgumentNullException.ThrowIfNull(url);
        var connectionFactory = new ConnectionFactory
        {
            Uri = new Uri(url),
            DispatchConsumersAsync = true,
        };

        return connectionFactory.CreateConnection();
    }
}